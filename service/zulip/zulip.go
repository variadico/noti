package zulip

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

var (
// API is the Pushbullet API endpoint.

)

type apiResponse struct {
	ID     int    `json:"id"`
	Code   string `json:"code"`
	Msg    string `json:"msg"`
	Result string `json:"result"`
	Stream string `json:"stream"`
}

// Notification is a pushbullet notification.
type Notification struct {
	Content  string
	Type     string
	To       string
	Endpoint string

	BotAPIKey       string
	BotEmailAddress string
	Client          *http.Client
}

// Send sends a Pushbullet notification.
func (n *Notification) Send() error {
	data := url.Values{}
	data.Set("type", n.Type)
	if n.Type != "stream" && n.Type != "private" {
		return fmt.Errorf("%s != stream || private", n.Type)
	}

	if n.Type == "stream" {
		data.Set("subject", "Noti")
	}

	data.Set("to", n.To)
	data.Set("content", n.Content)

	req, err := http.NewRequest("POST", n.Endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(n.BotEmailAddress, n.BotAPIKey)

	resp, err := n.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	if r.Result != "success" {
		return fmt.Errorf("zuliperror: %+v %+v", data, r)
	}

	return nil
}
