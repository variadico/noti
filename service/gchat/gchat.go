package gchat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Notification is a Google Chat notification.
type Notification struct {
	// AppURL is your Google Chat App's webhook URL
	AppURL string `json:"-"`

	// Text is the notification's message.
	Text string `json:"text"`

	Client *http.Client `json:"-"`
}

// Send triggers a Google Chat notification.
func (n *Notification) Send() error {
	if n.AppURL == "" {
		return errors.New("missing App URL")
	}
	if n.Text == "" {
		return errors.New("missing message text")
	}

	data := new(bytes.Buffer)
	if err := json.NewEncoder(data).Encode(n); err != nil {
		return err
	}

	resp, err := n.Client.Post(n.AppURL, "application/json", data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buff := new(bytes.Buffer)
	_, err = buff.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	s := buff.String()

	if s != "ok" {
		return fmt.Errorf("google chat api error: %s", s)
	}

	return nil
}
