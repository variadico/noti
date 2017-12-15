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
	"github.com/spf13/viper"
	"github.com/variadico/vbs"
)

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

	Root.Flags().SetInterspersed(false)

	Root.Flags().StringP("title", "t", "", "Notification title. Default is utility name.")
	Root.Flags().StringP("message", "m", "Done!", "Notification message. Default is 'Done!'.")

	Root.Flags().IntP("pwatch", "w", -1, "Trigger notification after PID disappears.")

	Root.Flags().BoolP("banner", "b", true, "Trigger a banner notification. Default is true. To disable this notification set this flag to false.")
	Root.Flags().BoolP("speech", "s", false, "Trigger a speech notification. Optionally, customize the voice with NOTI_VOICE.")
	Root.Flags().BoolP("hipchat", "i", false, "Trigger a HipChat notification. Requires NOTI_HIPCHAT_TOK and NOTI_HIPCHAT_DEST to be set.")
	Root.Flags().BoolP("pushbullet", "p", false, "Trigger a Pushbullet notification. Requires NOTI_PUSHBULLET_TOK to be set.")
	Root.Flags().BoolP("pushover", "o", false, "Trigger a Pushover notification. Requires NOTI_PUSHOVER_TOK and NOTI_PUSHOVER_DEST to be set.")
	Root.Flags().BoolP("pushsafer", "u", false, "Trigger a Pushsafer notification. Requires NOTI_PUSHSAFER_KEY to be set.")
	Root.Flags().BoolP("simplepush", "l", false, "Trigger a Simplepush notification. Requires NOTI_SIMPLEPUSH_KEY to be set. Optionally, customize ringtone and vibration with NOTI_SIMPLEPUSH_EVENT.")
	Root.Flags().BoolP("slack", "k", false, "Trigger a Slack notification. Requires NOTI_SLACK_TOK and NOTI_SLACK_DEST to be set.")
	Root.Flags().BoolP("bearychat", "c", false, "Trigger a BearyChat notification. Requries NOTI_BC_INCOMING_URI to be set.")

	Root.Flags().BoolP("version", "v", false, "Print noti version and exit.")
	Root.Flags().BoolP("help", "h", false, "Display help information and exit.")
	Root.Flags().BoolVar(&vbs.Enabled, "verbose", false, "Enable verbose mode.")
}

func rootMain(cmd *cobra.Command, args []string) error {
	v := viper.New()
	setNotiDefaults(v)
	bindNotiEnv(v)
	if err := setupConfigFile(v); err != nil {
		vbs.Println("Failed to read config file:", err)
	}
	v.BindPFlag("message", cmd.Flags().Lookup("message"))

	vbs.Println("Command:", args)
	if vbs.Enabled {
		printEnv()
	}

	if showVer, _ := cmd.Flags().GetBool("version"); showVer {
		fmt.Println(Version)
		checkForUpdates()
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

	if pid, _ := cmd.Flags().GetInt("pwatch"); pid != -1 {
		vbs.Println("Watching PID")
		err = pollPID(pid, 1*time.Second)
	} else {
		vbs.Println("Running command")
		err = runCommand(args, os.Stdin, os.Stdout, os.Stderr)
	}
	if err != nil {
		v.Set("message", err.Error())
		v.Set("nsuser.soundName", v.GetString("nsuser.soundNameFail"))
	}

	config := readEnv(os.Getenv("NOTI_DEFAULT"))
	err = readFlags(cmd.Flags(), config, os.Getenv("NOTI_DEFAULT") == "")
	if err != nil {
		return err
	}

	vbs.Println("Config:", config)
	vbs.Println("Viper:", v.AllSettings())
	var notis []notification
	message := v.GetString("message")

	if config["banner"] {
		notis = append(notis, getBanner(title, message, v))
	}

	if config["speech"] {
		notis = append(notis, getSpeech(title, message, v))
	}

	if config["bearychat"] {
		notis = append(notis, getBearyChat(title, message, v))
	}

	if config["hipchat"] {
		notis = append(notis, getHipChat(title, message, v))
	}

	if config["pushbullet"] {
		notis = append(notis, getPushbullet(title, message, v))
	}

	if config["pushover"] {
		notis = append(notis, getPushover(title, message, v))
	}

	if config["pushsafer"] {
		notis = append(notis, getPushsafer(title, message, v))
	}

	if config["simplepush"] {
		notis = append(notis, getSimplepush(title, message, v))
	}

	if config["slack"] {
		notis = append(notis, getSlack(title, message, v))
	}

	vbs.Println("Notifications:", len(notis))
	for _, n := range notis {
		if err := n.Send(); err != nil {
			log.Println(err)
		}
	}

	return nil
}

func checkForUpdates() error {
	// Draft releases and prereleases are not returned by this endpoint.
	const releaseAPI = "https://api.github.com/repos/variadico/noti/releases/latest"
	webClient := &http.Client{Timeout: 30 * time.Second}

	resp, err := webClient.Get(releaseAPI)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r struct {
		HTMLURL string `json:"html_url"`
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	fmt.Println("latest:", r.TagName)
	if r.TagName != Version {
		fmt.Println("download:", r.HTMLURL)
	}

	return nil
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

// printEnv prints all of the environment variables used in noti.
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
