package command

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/variadico/vbs"
)

// Draft releases and prereleases are not returned by this endpoint.
const githubReleasesURL = "https://api.github.com/repos/variadico/noti/releases/latest"

// notification is the interface for all notifications.
type notification interface {
	Send() error
}

// Root is the root noti command.
var Root = &cobra.Command{
	Use:  "noti [flags] [command [arguments]]",
	RunE: rootMain,

	SilenceErrors: true,
	SilenceUsage:  true,
}

// Version is the version of noti. This is set at compile time.
var Version string

func init() {
	var compiledManual string
	switch runtime.GOOS {
	case "darwin":
		compiledManual = fmt.Sprintf(manual, osxManual)
	case "linux", "freebsd":
		compiledManual = fmt.Sprintf(manual, linuxFreeBSDManual)
	}

	Root.SetUsageTemplate(compiledManual)
	defineFlags(Root.Flags())
}

func defineFlags(flags *pflag.FlagSet) {
	flags.SetInterspersed(false)

	flags.StringP("title", "t", "", "Notification title. Default is utility name.")
	flags.StringP("message", "m", "", "Notification message. Default is 'Done!'.")

	flags.IntP("pwatch", "w", -1, "Trigger notification after PID disappears.")

	flags.BoolP("banner", "b", false, "Trigger a banner notification. Default is true. To disable this notification set this flag to false.")
	flags.BoolP("speech", "s", false, "Trigger a speech notification. Optionally, customize the voice with NOTI_VOICE.")
	flags.BoolP("hipchat", "i", false, "Trigger a HipChat notification. Requires NOTI_HIPCHAT_TOK and NOTI_HIPCHAT_DEST to be set.")
	flags.BoolP("pushbullet", "p", false, "Trigger a Pushbullet notification. Requires NOTI_PUSHBULLET_TOK to be set.")
	flags.BoolP("pushover", "o", false, "Trigger a Pushover notification. Requires NOTI_PUSHOVER_TOK and NOTI_PUSHOVER_DEST to be set.")
	flags.BoolP("pushsafer", "u", false, "Trigger a Pushsafer notification. Requires NOTI_PUSHSAFER_KEY to be set.")
	flags.BoolP("simplepush", "l", false, "Trigger a Simplepush notification. Requires NOTI_SIMPLEPUSH_KEY to be set. Optionally, customize ringtone and vibration with NOTI_SIMPLEPUSH_EVENT.")
	flags.BoolP("slack", "k", false, "Trigger a Slack notification. Requires NOTI_SLACK_TOK and NOTI_SLACK_DEST to be set.")
	flags.BoolP("bearychat", "c", false, "Trigger a BearyChat notification. Requries NOTI_BC_INCOMING_URI to be set.")

	flags.BoolP("version", "v", false, "Print noti version and exit.")
	flags.BoolP("help", "h", false, "Display help information and exit.")
	flags.BoolVar(&vbs.Enabled, "verbose", false, "Enable verbose mode.")
}

func rootMain(cmd *cobra.Command, args []string) error {
	v := viper.New()
	if err := configureApp(v, cmd.Flags()); err != nil {
		vbs.Println("Config error:", err)
	}

	if vbs.Enabled {
		printEnv()
	}

	if showVer, _ := cmd.Flags().GetBool("version"); showVer {
		fmt.Println("noti version", Version)
		if latest, dl, err := latestRelease(githubReleasesURL); err != nil {
			vbs.Println("Failed get latest release:", err)
		} else if latest != Version {
			fmt.Println("Latest:", latest)
			fmt.Println("Download:", dl)
		}
		return nil
	}

	if showHelp, _ := cmd.Flags().GetBool("help"); showHelp {
		return cmd.Usage()
	}

	title, err := cmd.Flags().GetString("title")
	if err != nil {
		return err
	}
	if title == "" {
		title = commandName(args)
	}
	vbs.Println("title:", title)
	v.Set("title", title)

	if pid, _ := cmd.Flags().GetInt("pwatch"); pid != -1 {
		vbs.Println("Watching PID:", pid)
		err = pollPID(pid, 2*time.Second)
	} else {
		vbs.Println("Running command:", args)
		err = runCommand(args, os.Stdin, os.Stdout, os.Stderr)
	}
	if err != nil {
		v.Set("message", err.Error())
		v.Set("nsuser.soundName", v.GetString("nsuser.soundNameFail"))
	}

	enabled := enabledServices(v, cmd.Flags())
	vbs.Println("Services:", enabled)
	vbs.Println("Viper:", v.AllSettings())
	notis := getNotifications(v, enabled)

	vbs.Println(len(notis), "notifications queued")
	for _, n := range notis {
		if err := n.Send(); err != nil {
			log.Println(err)
		} else {
			vbs.Printf("Sent: %T\n", n)
		}
	}

	return nil
}

func latestRelease(u string) (string, string, error) {
	webClient := &http.Client{Timeout: 30 * time.Second}

	resp, err := webClient.Get(u)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var r struct {
		HTMLURL string `json:"html_url"`
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", "", err
	}

	return r.TagName, r.HTMLURL, nil
}

func commandName(args []string) string {
	switch len(args) {
	case 0:
		return "noti"
	case 1:
		return args[0]
	}

	if args[1][0] != '-' {
		// If the next arg isn't a flag, append a subcommand to the command
		// name.
		return fmt.Sprintf("%s %s", args[0], args[1])
	}

	return args[0]
}

func runCommand(args []string, sin io.Reader, sout, serr io.Writer) error {
	if len(args) == 0 {
		return nil
	}

	var cmd *exec.Cmd
	if _, err := exec.LookPath(args[0]); err != nil {
		// Maybe command is alias or builtin?
		cmd = subshellCommand(args)
		if cmd == nil {
			return err
		}
	} else {
		cmd = exec.Command(args[0], args[1:]...)
	}

	cmd.Stdin = sin
	cmd.Stdout = sout
	cmd.Stderr = serr
	return cmd.Run()
}

func subshellCommand(args []string) *exec.Cmd {
	shell := os.Getenv("SHELL")

	switch filepath.Base(shell) {
	case "bash", "zsh":
		args = append([]string{"-l", "-i", "-c"}, args...)
	default:
		return nil
	}

	return exec.Command(shell, args...)
}

func printEnv() {
	envs := []string{
		"NOTI_BC_INCOMING_URI",
		"NOTI_DEFAULT",
		"NOTI_HIPCHAT_DEST",
		"NOTI_HIPCHAT_TOK",
		"NOTI_PUSHBULLET_TOK",
		"NOTI_PUSHOVER_DEST",
		"NOTI_PUSHOVER_TOK",
		"NOTI_PUSHSAFER_KEY",
		"NOTI_SIMPLEPUSH_EVENT",
		"NOTI_SIMPLEPUSH_KEY",
		"NOTI_SLACK_DEST",
		"NOTI_SLACK_TOK",
		"NOTI_SOUND",
		"NOTI_SOUND_FAIL",
		"NOTI_VOICE",
	}

	for _, env := range envs {
		if val, set := os.LookupEnv(env); set {
			fmt.Printf("%s=%s\n", env, val)
		}
	}
}
