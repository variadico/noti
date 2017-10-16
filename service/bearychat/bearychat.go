package bearychat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type apiResponse struct {
	Code  int    `json:"code"`
	Error string `json:"string"`
}

// Notification is a BearyChat notification.
type Notification struct {
	Text         string `json:"text"`
	Notification string `json:"notification,omitempty"`
	Markdown     bool   `json:"markdown,omitempty"`
	Channel      string `json:"channel,omitempty"`
	User         string `json:"user,omitempty"`

	IncomingHookURI string       `json:"-"`
	Client          *http.Client `json:"-"`
}

// Send sends a message request to BearyChat's incoming hook.
func (n *Notification) Send() error {
	if n.Text == "" {
		return errors.New("missing text")
	}

	payload := new(bytes.Buffer)
	if err := json.NewEncoder(payload).Encode(n); err != nil {
		return err
	}

	resp, err := n.Client.Post(n.IncomingHookURI, "application/json", payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return fmt.Errorf("decoding response: %s", err)
	}

	if r.Code != 0 {
		return errors.New(r.Error)
	}

	return nil
}

// SetMessage sets a notification's message.
func (n *Notification) SetMessage(m string) {
	n.Text = m
}

// GetMessage gets a notification's message.
func (n *Notification) GetMessage() string {
	return n.Text
}
