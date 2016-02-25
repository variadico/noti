package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	version = "2.1.0"

	defaultEnv       = "NOTI_DEFAULT"
	hipChatDestEnv   = "NOTI_HIPCHAT_DEST"
	hipChatTokEnv    = "NOTI_HIPCHAT_TOK"
	pushbulletTokEnv = "NOTI_PUSHBULLET_TOK"
	pushoverTokEnv   = "NOTI_PUSHOVER_TOK"
	pushoverDestEnv  = "NOTI_PUSHOVER_DEST"
	slackDestEnv     = "NOTI_SLACK_DEST"
	slackTokEnv      = "NOTI_SLACK_TOK"
	soundEnv         = "NOTI_SOUND"
	soundFailEnv     = "NOTI_SOUND_FAIL"
	voiceEnv         = "NOTI_VOICE"

	webTimeout = 30 * time.Second
)

var (
	title       = flag.String("t", "noti", "")
	message     = flag.String("m", "Done!", "")
	showVersion = flag.Bool("v", false, "")
	showHelp    = flag.Bool("h", false, "")

	// Notifications
	banner     = flag.Bool("b", false, "")
	hipChat    = flag.Bool("i", false, "")
	pushbullet = flag.Bool("p", false, "")
	pushover   = flag.Bool("o", false, "")
	slack      = flag.Bool("k", false, "")
	speech     = flag.Bool("s", false, "")

	utilityFailed bool

	webClient = &http.Client{Timeout: webTimeout}
)

func init() {
	flag.StringVar(title, "title", "noti", "")
	flag.StringVar(message, "message", "Done!", "")
	flag.BoolVar(showVersion, "version", false, "")
	flag.BoolVar(showHelp, "help", false, "")

	// Notifications
	flag.BoolVar(banner, "banner", false, "")
	flag.BoolVar(hipChat, "hipchat", false, "")
	flag.BoolVar(speech, "speech", false, "")
	flag.BoolVar(pushbullet, "pushbullet", false, "")
	flag.BoolVar(pushover, "pushover", false, "")
	flag.BoolVar(slack, "slack", false, "")
}

func main() {
	log.SetFlags(0)
	flag.Parse()

	if *showVersion {
		fmt.Println("noti version", version)
		return
	}
	if *showHelp {
		flag.Usage()
		return
	}

	runUtility()

	if defs := strings.TrimSpace(os.Getenv(defaultEnv)); defs != "" {
		*banner = strings.Contains(defs, "banner")
		*hipChat = strings.Contains(defs, "hipchat")
		*pushbullet = strings.Contains(defs, "pushbullet")
		*pushover = strings.Contains(defs, "pushover")
		*speech = strings.Contains(defs, "speech")
		*slack = strings.Contains(defs, "slack")
	} else {
		var strVal string
		var explicitSet bool

		if userSet("b") {
			strVal = flag.Lookup("b").Value.String()
			explicitSet = true
		} else if userSet("banner") {
			strVal = flag.Lookup("banner").Value.String()
			explicitSet = true
		}

		// If the user explicitly set -banner, then use the value that the user
		// set, but if no banner flag was set, then the default is true.
		if explicitSet {
			// Ignoring error, false on error is fine.
			*banner, _ = strconv.ParseBool(strVal)
		} else {
			*banner = true
		}
	}

	ns := []struct {
		run bool
		fn  func() error
	}{
		{*banner, bannerNotify},
		{*hipChat, hipChatNotify},
		{*pushbullet, pushbulletNotify},
		{*pushover, pushoverNotify},
		{*speech, speechNotify},
		{*slack, slackNotify},
	}

	for _, n := range ns {
		if !n.run {
			continue
		}

		if err := n.fn(); err != nil {
			log.Println(err)
		}
	}
}

func runUtility() {
	var cmd *exec.Cmd

	args := flag.Args()
	if len(args) < 1 {
		return
	}

	cmd = exec.Command(args[0], args[1:]...)
	*title = args[0]

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if exerr, is := err.(*exec.ExitError); is {
		if !exerr.Success() {
			*title = *title + " failed"
		}
	}
	if err != nil {
		utilityFailed = true
		*message = err.Error()
	}
}

func userSet(target string) bool {
	var set bool

	flag.Visit(func(f *flag.Flag) {
		if f.Name == target {
			set = true
		}
	})

	return set
}
