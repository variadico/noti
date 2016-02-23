package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func pushbulletNotify() {
	accessToken := os.Getenv(pushbulletTokEnv)
	if accessToken == "" {
		log.Fatalf("Missing access token, %s must be set", pushbulletTokEnv)
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

	if _, err = webClient.Do(req); err != nil {
		log.Fatal(err)
	}
}

func slackNotify() {
	accessToken := os.Getenv(slackTokEnv)
	if accessToken == "" {
		log.Fatalf("Missing access token, %s must be set", slackTokEnv)
	}

	dest := os.Getenv(slackDestEnv)
	if dest == "" {
		log.Fatalf("Missing destination, %s must be set", slackDestEnv)
	}

	vals := make(url.Values)
	vals.Set("token", accessToken)
	vals.Set("text", fmt.Sprintf("%s\n%s", *title, *message))
	vals.Set("username", "noti")
	vals.Set("channel", dest)

	resp, err := webClient.PostForm("https://slack.com/api/chat.postMessage", vals)
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

func hipChatNotify() {
	accessToken := os.Getenv(hipChatTokEnv)
	if accessToken == "" {
		log.Fatalf("Missing access token, %s must be set", hipChatTokEnv)
	}

	dest := os.Getenv(hipChatDestEnv)
	if dest == "" {
		log.Fatalf("Missing destination, %s must be set", hipChatDestEnv)
	}

	payload := bytes.NewBuffer([]byte(fmt.Sprintf(
		`{"message":%q,"message_format":"text"}`,
		fmt.Sprintf("%s\n%s", *title, *message),
	)))

	ep := fmt.Sprintf("https://api.hipchat.com/v2/room/%s/notification", dest)
	req, err := http.NewRequest("POST", ep, payload)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := webClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	r := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&r); err == io.EOF {
		return
	} else if err != nil {
		resp.Body.Close()
		log.Fatal(err)
	}
	resp.Body.Close()

	if err, exists := r["error"]; exists {
		if m, is := err.(map[string]interface{}); is {
			log.Fatal("HipChat API error: ", m["message"])
		}
	}

}
