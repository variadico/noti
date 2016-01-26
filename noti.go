package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	pushbulletEnv   = "NOTI_PUSHBULLET_TOK"
	slackEnv        = "NOTI_SLACK_TOK"
	slackChannelEnv = "NOTI_SLACK_CHAN"
	voiceEnv        = "NOTI_VOICE"
	soundEnv        = "NOTI_SOUND"
	defaultEnv      = "NOTI_DEFAULT"

	version = "v2dev"
)

var (
	title       = flag.String("t", "noti", "")
	message     = flag.String("m", "Done!", "")
	showVersion = flag.Bool("v", false, "")
	showHelp    = flag.Bool("h", false, "")

	// Notifications
	pushbullet = flag.Bool("p", false, "")
	speech     = flag.Bool("s", false, "")
	slack      = flag.Bool("S", false, "")
)

func init() {
	flag.StringVar(title, "title", "noti", "")
	flag.StringVar(message, "message", "Done!", "")
	flag.BoolVar(showVersion, "version", false, "")
	flag.BoolVar(showHelp, "help", false, "")

	// Notifications
	flag.BoolVar(speech, "speech", false, "")
	flag.BoolVar(pushbullet, "pushbullet", false, "")
	flag.BoolVar(slack, "slack", false, "")
}

func main() {
	log.SetFlags(0)
	flag.Parse()

	if *showVersion {
		fmt.Printf("noti version %s\n", version)
		return
	}
	if *showHelp {
		flag.Usage()
		return
	}

	switch strings.ToLower(os.Getenv(defaultEnv)) {
	case "slack":
		slackNotify()
		return
	case "pushbullet":
		pushbulletNotify()
		return
	case "speech":
		speechNotify()
		return
	case "desktop":
		desktopNotify()
		return
	}

	switch {
	case *slack:
		slackNotify()
	case *pushbullet:
		pushbulletNotify()
	case *speech:
		speechNotify()
	default:
		desktopNotify()
	}
}

func runUtility() {
	var cmd *exec.Cmd

	if args := flag.Args(); len(args) < 1 {
		return
	} else {
		cmd = exec.Command(args[0], args[1:]...)
		*title = args[0]
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		*title = *title + " failed"
		*message = err.Error()
	}
}
