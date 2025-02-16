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
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Draft releases and prereleases are not returned by this endpoint.
const githubReleasesURL = "https://api.github.com/repos/variadico/noti/releases/latest"

// notification is the interface for all notifications.
type notification interface {
	Send() error
}

// Root is the root noti command.
var Root = &cobra.Command{
	Long:    "noti - Monitor a process and trigger a notification",
	Use:     "noti [flags] [utility [args...]]",
	Example: "noti tar -cjf music.tar.bz2 Music/\nclang foo.c; noti",
	RunE:    rootMain,

	SilenceErrors: true,
	SilenceUsage:  true,
}

// Version is the version of noti. This is set at compile time with Make.
var Version string

// InitFlags initializes Root's command line flags.
func InitFlags(flags *pflag.FlagSet) {
	flags.SetInterspersed(false)
	flags.SortFlags = false

	flags.StringP("title", "t", "", "Set notification title. Default is utility name.")
	flags.StringP("message", "m", "", `Set notification message. Default is "Done!". Read from stdin with "-".`)
	flags.BoolP("time", "e", false, "Show execution time in message.")

	flags.BoolP("banner", "b", false, "Trigger a banner notification. This is enabled by default.")
	flags.BoolP("speech", "s", false, "Trigger a speech notification.")
	flags.BoolP("bearychat", "c", false, "Trigger a BearyChat notification.")
	flags.Bool("keybase", false, "Trigger a Keybase notification.")
	flags.BoolP("pushbullet", "p", false, "Trigger a Pushbullet notification.")
	flags.BoolP("pushover", "o", false, "Trigger a Pushover notification.")
	flags.BoolP("pushsafer", "u", false, "Trigger a Pushsafer notification.")
	flags.BoolP("simplepush", "l", false, "Trigger a Simplepush notification.")
	flags.BoolP("slack", "k", false, "Trigger a Slack notification.")
	flags.BoolP("gchat", "n", false, "Trigger a Google Chat notification.")
	flags.BoolP("chanify", "i", false, "Trigger a Chanify notification.")
	flags.BoolP("mattermost", "a", false, "Trigger a Mattermost notification")
	flags.BoolP("telegram", "g", false, "Trigger a Telegram notification")
	flags.BoolP("zulip", "z", false, "Trigger a Zulip notification")
	flags.Bool("twilio", false, "Trigger a twilio SMS notification")
	flags.Bool("ntfy", false, "Trigger a Ntfy notification")
	flags.Bool("bark", false, "Trigger a Bark notification")

	flags.IntP("pwatch", "w", -1, "Monitor a process by PID and trigger a notification when the pid disappears.")
	flags.StringP("file", "f", "", "Path to noti.yaml configuration file.")
	flags.BoolVar(&vbsEnabled, "verbose", false, "Enable verbose mode.")
	flags.BoolP("version", "v", false, "Print noti version and exit.")
	flags.BoolP("help", "h", false, "Print noti help and exit.")
}

func rootMain(cmd *cobra.Command, args []string) error {
	vbsPrintln("os.Args:", os.Args)

	v := viper.New()
	if err := configureApp(v, cmd.Flags()); err != nil {
		vbsPrintln("Failed to configure:", err)
	}

	if vbsEnabled {
		printEnv()
	}

	if showVer, _ := cmd.Flags().GetBool("version"); showVer {
		fmt.Println("noti version", Version)
		if latest, dl, err := latestRelease(githubReleasesURL); err != nil {
			vbsPrintln("Failed get latest release:", err)
		} else if latest != Version {
			fmt.Println("Latest:", latest)
			fmt.Println("Download:", dl)
		}
		return nil
	}

	if showHelp, _ := cmd.Flags().GetBool("help"); showHelp {
		return cmd.Help()
	}

	title, err := cmd.Flags().GetString("title")
	if err != nil {
		return err
	}
	if title == "" {
		vbsPrintln("Title from flags is empty, getting title from command name")
		title = commandName(args)
	}
	v.Set("title", title)

	var (
		cmdErr  error
		cmdTime time.Duration
	)
	if pid, _ := cmd.Flags().GetInt("pwatch"); pid != -1 {
		vbsPrintln("Watching PID:", pid)
		if err := pollPID(pid, 2*time.Second); err != nil {
			return err
		}
	} else if msg, _ := cmd.Flags().GetString("message"); msg == "-" {
		buf, err := io.ReadAll(os.Stdin)
		if err != nil {
			// buffer overflow
			return err
		}
		v.Set("message", string(buf))
	} else {
		vbsPrintln("Running command:", args)
		timeBefore := time.Now()
		cmdErr = runCommand(args, os.Stdin, os.Stdout, os.Stderr)
		cmdTime = time.Since(timeBefore).Round(time.Second)
	}
	if cmdErr != nil {
		v.Set("message", cmdErr.Error())
		v.Set("nsuser.soundName", v.GetString("nsuser.soundNameFail"))
	}
	if enabledTime(v, cmd.Flags()) {
		v.Set("message", fmt.Sprintf("%s (%s)", v.GetString("message"), cmdTime))
	}

	vbsPrintln("Title:", v.GetString("title"))
	vbsPrintln("Message:", v.GetString("message"))
	vbsPrintln("Time:", enabledTime(v, cmd.Flags()))

	enabled := enabledServices(v, cmd.Flags())
	vbsPrintln("Services:", enabled)
	vbsPrintln("Viper:", v.AllSettings())
	notis := getNotifications(v, enabled)

	vbsPrintln(len(notis), "notifications queued")
	for _, n := range notis {
		if err := n.Send(); err != nil {
			log.Println(err)
		} else {
			vbsPrintf("Sent: %T\n", n)
		}
	}

	return cmdErr
}

func enabledTime(v *viper.Viper, flags *pflag.FlagSet) bool {
	if measureTime, _ := flags.GetBool("time"); measureTime {
		return true
	}
	if v.GetBool("time") {
		return true
	}
	return false
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

	missing := make([]string, 0, 2)
	if r.HTMLURL == "" {
		missing = append(missing, "HTML URL")
	}
	if r.TagName == "" {
		missing = append(missing, "tag name")
	}
	if len(missing) > 0 {
		return "", "", fmt.Errorf("latest release: missing %q", strings.Join(missing, ","))
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
		cmd = shellCommand(args)
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

// shellCommand returns a shell alias or builtin command, as opposed to a
// program installed on the filesystem. This is needed to allow people to use
// noti with aliases or builtins.
func shellCommand(args []string) *exec.Cmd {
	shell := os.Getenv("SHELL")

	switch filepath.Base(shell) {
	case "bash", "zsh":
		args = append([]string{"-i", "-c"}, args...)
	default:
		return nil
	}

	return exec.Command(shell, args...)
}

func printEnv() {
	alloc := len(keyEnvBindings) + len(keyEnvBindingsDeprecated)
	envs := make([]string, 0, alloc)
	for _, e := range keyEnvBindings {
		envs = append(envs, e)
	}
	for _, e := range keyEnvBindingsDeprecated {
		envs = append(envs, e)
	}

	for _, env := range envs {
		if val, set := os.LookupEnv(env); set {
			fmt.Printf("%s=%s\n", env, val)
		}
	}
}
