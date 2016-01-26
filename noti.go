package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

const (
	pushbulletEnv = "NOTI_PUSHBULLET_TOK"
	slackEnv      = "NOTI_SLACK_TOK"
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

func pushbulletNotify() {
	runUtility()

	accessToken := os.Getenv(pushbulletEnv)
	if accessToken == "" {
		log.Fatalf("Missing access token, %s must be set", pushbulletEnv)
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

func slackNotify() {
	runUtility()

	accessToken := os.Getenv(slackEnv)
	if accessToken == "" {
		log.Fatalf("Missing access token, %s must be set", slackEnv)
	}

	vals := make(url.Values)
	vals.Set("token", accessToken)
	vals.Set("text", fmt.Sprintf("%s\n%s", *title, *message))
	vals.Set("username", "noti")
	vals.Set("channel", "#random")

	resp, err := http.PostForm("https://slack.com/api/chat.postMessage", vals)
	if err != nil {
		log.Fatal(err)
	}

	r := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		resp.Body.Close()
		log.Fatal(err)
	}
	resp.Body.Close()

	if r["ok"] == false {
		log.Fatal("Slack API error: ", r["error"])
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
