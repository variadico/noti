package command

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/variadico/noti/service"
	"github.com/variadico/vbs"
)

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
	Root.Flags().StringP("message", "m", "", "Notification message. Default is 'Done!'.")

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

	message, err := cmd.Flags().GetString("message")
	if err != nil {
		return err
	}
	if message == "" {
		message = "Done!"
	}

	pass, fail := macOSSounds(os.Getenv("NOTI_SOUND"), os.Getenv("NOTI_SOUND_FAIL"))
	sound := pass

	if pid, _ := cmd.Flags().GetInt("pwatch"); pid != -1 {
		vbs.Println("Watching PID")
		err = pollPID(pid, 1*time.Second)
	} else {
		vbs.Println("Running command")
		err = runCommand(args)
	}
	if err != nil {
		message = err.Error()
		sound = fail
	}

	config := readEnv(os.Getenv("NOTI_DEFAULT"))
	err = readFlags(cmd.Flags(), config, os.Getenv("NOTI_DEFAULT") == "")
	if err != nil {
		return err
	}

	vbs.Println("Config:", config)
	var notis []service.Notification

	if config["banner"] {
		notis = append(notis, getBanner(title, message, sound))
	}

	if config["speech"] {
		notis = append(notis, getSpeech(title, message, os.Getenv("NOTI_VOICE")))
	}

	if config["bearychat"] {
		notis = append(notis, getBearyChat(title, message,
			os.Getenv("NOTI_BC_INCOMING_URI")))
	}

	if config["hipchat"] {
		notis = append(notis, getHipChat(title, message,
			os.Getenv("NOTI_HIPCHAT_TOK"), os.Getenv("NOTI_HIPCHAT_DEST")))
	}

	if config["pushbullet"] {
		notis = append(notis, getPushbullet(title, message,
			os.Getenv("NOTI_PUSHBULLET_TOK")))
	}

	if config["pushover"] {
		notis = append(notis, getPushover(title, message,
			os.Getenv("NOTI_PUSHOVER_TOK"), os.Getenv("NOTI_PUSHOVER_DEST")))
	}

	if config["pushsafer"] {
		notis = append(notis, getPushsafer(title, message,
			os.Getenv("NOTI_PUSHSAFER_KEY")))
	}

	if config["simplepush"] {
		notis = append(notis, getSimplepush(title, message,
			os.Getenv("NOTI_SIMPLEPUSH_KEY"), os.Getenv("NOTI_SIMPLEPUSH_EVENT")))
	}

	if config["slack"] {
		notis = append(notis, getSlack(title, message,
			os.Getenv("NOTI_SLACK_TOK"), os.Getenv("NOTI_SLACK_DEST")))
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

func runCommand(args []string) error {
	if len(args) == 0 {
		return nil
	}

	if _, err := exec.LookPath(args[0]); err != nil {
		exp, expErr := expandAlias(args[0])
		if expErr != nil {
			return err
		}

		args = append(exp, args[1:]...)
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func expandAlias(alias string) ([]string, error) {
	shell := os.Getenv("SHELL")

	cmd := exec.Command(shell, "-l", "-i", "-c", "which "+alias)
	expanded, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	exp := strings.TrimSpace(string(expanded))
	trimLen := fmt.Sprintf("%s: aliased to ", alias)
	exp = exp[len(trimLen):]

	return strings.Split(exp, " "), nil
}
