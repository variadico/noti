package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	pushbulletEnv = "NOTI_PB_ACCESS"
	voiceEnv      = "NOTI_VOICE"
	soundEnv      = "NOTI_SOUND"
	defaultEnv    = "NOTI_DEFAULT"

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
)

func init() {
	flag.StringVar(title, "title", "noti", "")
	flag.StringVar(message, "message", "Done!", "")
	flag.BoolVar(showVersion, "version", false, "")
	flag.BoolVar(showHelp, "help", false, "")

	// Notifications
	flag.BoolVar(speech, "speech", false, "")
	flag.BoolVar(pushbullet, "pushbullet", false, "")
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
	case *pushbullet:
		pushbulletNotify()
	case *speech:
		speechNotify()
	default:
		desktopNotify()
	}
}

func pushbulletNotify() {
	runUtility()

	accessToken := os.Getenv(pushbulletEnv)
	if accessToken == "" {
		log.Fatal("Missing Pushbullet access token, NOTI_PB_ACCESS must be set")
	}

	payload := bytes.NewBuffer([]byte(fmt.Sprintf(
		`{"body":%q,"title":%q,"type":"note"}`, *message, *title,
	)))

	req, err := http.NewRequest("POST", "https://api.pushbullet.com/v2/pushes", payload)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Access-Token", accessToken)
	req.Header.Set("Content-Type", "application/json")

	if _, err = http.DefaultClient.Do(req); err != nil {
		log.Fatal(err)
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
