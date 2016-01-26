package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

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

	if ch := os.Getenv(slackChannelEnv); ch == "" {
		vals.Set("channel", "#random")
	} else {
		vals.Set("channel", ch)
	}

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
