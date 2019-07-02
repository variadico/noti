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
	Long:    "noti - Monitor a process and trigger a notification",
	Use:     "noti [flags] [utility [args...]]",
	Example: "noti tar -cjf music.tar.bz2 Music/\nclang foo.c; noti",
	RunE:    rootMain,

	SilenceErrors: true,
	SilenceUsage:  true,
}

// Version is the version of noti. This is set at compile time with Make.
var Version string

func init() {
	defineFlags(Root.Flags())
}

func defineFlags(flags *pflag.FlagSet) {
	flags.SetInterspersed(false)
	flags.SortFlags = false

	flags.StringP("title", "t", "", "Set notification title. Default is utility name.")
	flags.StringP("message", "m", "", `Set notification message. Default is "Done!".`)

	flags.BoolP("banner", "b", false, "Trigger a banner notification. This is enabled by default.")
	flags.BoolP("speech", "s", false, "Trigger a speech notification.")
	flags.BoolP("bearychat", "c", false, "Trigger a BearyChat notification.")
	flags.BoolP("hipchat", "i", false, "Trigger a HipChat notification.")
	flags.BoolP("pushbullet", "p", false, "Trigger a Pushbullet notification.")
	flags.BoolP("pushover", "o", false, "Trigger a Pushover notification.")
	flags.BoolP("pushsafer", "u", false, "Trigger a Pushsafer notification.")
	flags.BoolP("simplepush", "l", false, "Trigger a Simplepush notification.")
	flags.BoolP("slack", "k", false, "Trigger a Slack notification.")
	flags.BoolP("mattermost", "a", false, "Trigger a Mattermost notification")
	flags.BoolP("telegram", "g", false, "Trigger a Telegram notification")

	flags.IntP("pwatch", "w", -1, "Monitor a process by PID and trigger a notification when the pid disappears.")

	flags.StringP("file", "f", "", "Path to noti.yaml configuration file.")
	flags.BoolVar(&vbs.Enabled, "verbose", false, "Enable verbose mode.")
	flags.BoolP("version", "v", false, "Print noti version and exit.")
	flags.BoolP("help", "h", false, "Print noti help and exit.")
}

func rootMain(cmd *cobra.Command, args []string) error {
	vbs.Println("os.Args:", os.Args)

	v := viper.New()
	if err := configureApp(v, cmd.Flags()); err != nil {
		vbs.Println("Failed to configure:", err)
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
		return cmd.Help()
	}

	title, err := cmd.Flags().GetString("title")
	if err != nil {
		return err
	}
	if title == "" {
		vbs.Println("Title from flags is empty, getting title from command name")
		title = commandName(args)
	}
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

	vbs.Println("Title:", v.GetString("title"))
	vbs.Println("Message:", v.GetString("message"))

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
	var envs []string
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
