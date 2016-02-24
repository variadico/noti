package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	pushbulletAPI = "https://api.pushbullet.com/v2/pushes"
	slackAPI      = "https://slack.com/api/chat.postMessage"
	hipChatAPI    = "https://api.hipchat.com/v2/room/%s/notification"
	pushoverAPI   = "https://api.pushover.net/1/messages.json"
)

func pushbulletNotify() error {
	accessToken := os.Getenv(pushbulletTokEnv)
	if accessToken == "" {
		return fmt.Errorf("Missing access token, %s must be set", pushbulletTokEnv)
	}

	payload := bytes.NewBuffer([]byte(fmt.Sprintf(
		`{"body":%q,"title":%q,"type":"note"}`, *message, *title,
	)))

	req, err := http.NewRequest("POST", pushbulletAPI, payload)
	if err != nil {
		return err
	}
	req.Header.Set("Access-Token", accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := webClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func slackNotify() error {
	accessToken := os.Getenv(slackTokEnv)
	if accessToken == "" {
		return fmt.Errorf("Missing access token, %s must be set", slackTokEnv)
	}

	dest := os.Getenv(slackDestEnv)
	if dest == "" {
		return fmt.Errorf("Missing destination, %s must be set", slackDestEnv)
	}

	vals := make(url.Values)
	vals.Set("token", accessToken)
	vals.Set("text", fmt.Sprintf("%s\n%s", *title, *message))
	vals.Set("username", "noti")
	vals.Set("channel", dest)

	resp, err := webClient.PostForm(slackAPI, vals)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r struct {
		OK    bool
		Error string
	}

	if err := json.NewDecoder(resp.Body).Decode(&r); err == io.EOF {
		return nil
	} else if err != nil {
		return err
	}

	if !r.OK {
		return fmt.Errorf("Slack API error: %s", r.Error)
	}

	return nil
}

func hipChatNotify() error {
	accessToken := os.Getenv(hipChatTokEnv)
	if accessToken == "" {
		return fmt.Errorf("Missing access token, %s must be set", hipChatTokEnv)
	}

	dest := os.Getenv(hipChatDestEnv)
	if dest == "" {
		return fmt.Errorf("Missing destination, %s must be set", hipChatDestEnv)
	}

	payload := bytes.NewBuffer([]byte(fmt.Sprintf(
		`{"message":%q,"message_format":"text"}`,
		fmt.Sprintf("%s\n%s", *title, *message),
	)))

	req, err := http.NewRequest("POST", fmt.Sprintf(hipChatAPI, dest), payload)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := webClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r struct {
		Error struct {
			Code    int
			Message string
			Type    string
		}
	}

	if err := json.NewDecoder(resp.Body).Decode(&r); err == io.EOF {
		return nil
	} else if err != nil {
		return err
	}

	if m := r.Error.Message; m != "" {
		return fmt.Errorf("HipChat API error: %s", m)
	}

	return nil
}

func pushoverNotify() error {
	accessToken := os.Getenv(pushoverTokEnv)
	if accessToken == "" {
		return fmt.Errorf("Missing access token, %s must be set", pushoverTokEnv)
	}

	dest := os.Getenv(pushoverDestEnv)
	if dest == "" {
		return fmt.Errorf("Missing destination, %s must be set", pushoverDestEnv)
	}

	vals := make(url.Values)
	vals.Set("token", accessToken)
	vals.Set("user", dest)
	vals.Set("message", *message)
	vals.Set("title", *title)

	resp, err := webClient.PostForm(pushoverAPI, vals)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r struct {
		Errors  []string
		Info    string
		Request string
		Status  int
		Token   string
	}

	if err := json.NewDecoder(resp.Body).Decode(&r); err == io.EOF {
		return nil
	} else if err != nil {
		return err
	}

	if r.Status != 1 {
		return fmt.Errorf("Pushover API error: %s", strings.Join(r.Errors, ": "))
	} else if strings.Contains(r.Info, "no active devices") {
		return fmt.Errorf("Pushover API error: %s", r.Info)
	}

	return nil
}
