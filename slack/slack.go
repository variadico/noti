package slack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/variadico/noti"
)

const (
	tokEnv  = "NOTI_SLACK_TOK"
	destEnv = "NOTI_SLACK_DEST"

	API = "https://slack.com/api/chat.postMessage"
)

var (
	errNoTok  = noti.ConfigErrror{Env: tokEnv, Reason: "missing"}
	errNoDest = noti.ConfigErrror{Env: destEnv, Reason: "missing"}
)

type configuration struct {
	accessToken string
	destination string
}

type apiResponse struct {
	OK        bool   `json:"ok"`
	Channel   string `json:"channel"`
	Timestamp string `json:"ts"`
	Message   struct {
		Text     string `json:"text"`
		Username string `json:"username"`
		Icons    struct {
			Emoji   string `json:"emoji"`
			Image64 string `json:"image_64"`
		} `json:"icons"`
		Type      string `json:"type"`
		Subtype   string `json:"subtype"`
		Timestamp string `json:"ts"`
	} `json:"message"`
	Error string `json:"error"`
}

func envConfig(env noti.EnvGetter) (configuration, error) {
	tok := env.Get(tokEnv)
	if tok == "" {
		return configuration{}, errNoTok
	}

	dest := env.Get(destEnv)
	if dest == "" {
		return configuration{}, errNoDest
	}

	return configuration{
		accessToken: tok,
		destination: dest,
	}, nil
}

func Notify(n noti.Notification) error {
	config, err := envConfig(n.Config)
	if err != nil {
		return err
	}

	vals := make(url.Values)
	vals.Set("channel", config.destination)
	vals.Set("icon_emoji", ":rocket:")
	vals.Set("text", fmt.Sprintf("%s\n%s", n.Title, n.Message))
	vals.Set("token", config.accessToken)
	vals.Set("username", "noti")

	webClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := webClient.PostForm(n.API, vals)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return fmt.Errorf("decoding response: %s", err)
	}

	if !r.OK {
		return noti.APIError{Site: "Slack", Msg: r.Error}
	}

	return nil
}
