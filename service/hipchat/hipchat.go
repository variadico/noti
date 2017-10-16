package hipchat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// API is the HipChat API endpoint.
var API = "https://api.hipchat.com/v2/room/%s/notification"

type apiResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

// Notification is a HipChat notification.
type Notification struct {
	Message       string `json:"message"`
	MessageFormat string `json:"message_format"`

	Token       string       `json:"-"`
	Destination string       `json:"-"`
	Client      *http.Client `json:"-"`
}

// Send triggers a HipChat notification.
func (n *Notification) Send() error {
	payload := new(bytes.Buffer)
	if err := json.NewEncoder(payload).Encode(n); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf(API, n.Destination), payload)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", n.Token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err == io.EOF {
		return nil
	} else if err != nil {
		return err
	}

	if m := r.Error.Message; m != "" {
		return errors.New(m)
	}

	return nil
}

// SetMessage sets a notification's message.
func (n *Notification) SetMessage(m string) {
	n.Message = m
}

// GetMessage gets a notification's message.
func (n *Notification) GetMessage() string {
	return n.Message
}
