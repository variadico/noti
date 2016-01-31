package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	version = "2.0.0"

	defaultEnv      = "NOTI_DEFAULT"
	pushbulletEnv   = "NOTI_PUSHBULLET_TOK"
	slackChannelEnv = "NOTI_SLACK_DEST"
	slackEnv        = "NOTI_SLACK_TOK"
	soundEnv        = "NOTI_SOUND"
	voiceEnv        = "NOTI_VOICE"
)

var (
	title       = flag.String("t", "noti", "")
	message     = flag.String("m", "Done!", "")
	showVersion = flag.Bool("v", false, "")
	showHelp    = flag.Bool("h", false, "")

	// Notifications
	banner     = flag.Bool("b", false, "")
	pushbullet = flag.Bool("p", false, "")
	speech     = flag.Bool("s", false, "")
	slack      = flag.Bool("k", false, "")
)

func init() {
	flag.StringVar(title, "title", "noti", "")
	flag.StringVar(message, "message", "Done!", "")
	flag.BoolVar(showVersion, "version", false, "")
	flag.BoolVar(showHelp, "help", false, "")

	// Notifications
	flag.BoolVar(banner, "banner", false, "")
	flag.BoolVar(speech, "speech", false, "")
	flag.BoolVar(pushbullet, "pushbullet", false, "")
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
		*speech = strings.Contains(defs, "speech")
		*pushbullet = strings.Contains(defs, "pushbullet")
		*slack = strings.Contains(defs, "slack")
	} else {
		var explicitSet bool
		var val bool

		flag.Visit(func(f *flag.Flag) {
			if f.Name == "b" || f.Name == "banner" {
				explicitSet = true
				// Ignoring error, false on error is fine.
				val, _ = strconv.ParseBool(f.Value.String())
			}
		})

		// If the user explicitly set -banner, then use the value that the user
		// set, but if no banner flag was set, then the default is true.
		if explicitSet {
			*banner = val
		} else {
			*banner = true
		}
	}

	if *banner {
		bannerNotify()
	}
	if *speech {
		speechNotify()
	}
	if *pushbullet {
		pushbulletNotify()
	}
	if *slack {
		slackNotify()
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
		*message = err.Error()
	}
}
