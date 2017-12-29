package pushover

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

var (
	// API is the Pushover API endpoint.
	API = "https://api.pushover.net/1/messages.json"
)

type apiResponse struct {
	Info    string   `json:"info"`
	Status  int      `json:"status"`
	Request string   `json:"request"`
	Errors  []string `json:"errors"`
	Token   string   `json:"token"`
}

// Notification is a pushover notification.
type Notification struct {
	Message  string
	Title    string
	APIToken string
	UserKey  string

	Client *http.Client
}

// Send sends a pushover notification.
func (n *Notification) Send() error {
	vals := make(url.Values)
	vals.Set("token", n.APIToken)
	vals.Set("user", n.UserKey)
	vals.Set("message", n.Message)
	vals.Set("title", n.Title)

	resp, err := n.Client.PostForm(API, vals)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	if r.Status != 1 {
		return errors.New(strings.Join(r.Errors, ": "))
	} else if strings.Contains(r.Info, "no active devices") {
		return errors.New(r.Info)
	}

	return nil
}
