package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type apiResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		MessageID int `json:"message_id"`
		Chat      struct {
			ID       int64  `json:"id"`
			Title    string `json:"title"`
			Username string `json:"username"`
			Type     string `json:"type"`
		} `json:"chat"`
		Date int64  `json:"date"`
		Text string `json:"text"`
	} `json:"result"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
}

// Notification is a Telegram notification.
type Notification struct {
	ChatID  string       `json:"chat_id"`
	Message string       `json:"text"`
	Token   string       `json:"-"`
	Client  *http.Client `json:"-"`
}

// API is the Telegram API endpoint.
var API = "https://api.telegram.org"

// Send sends a Telegram notification.
func (n *Notification) Send() error {
	if n.ChatID == "" {
		return errors.New("missing chat id")
	}

	if n.Token == "" {
		return errors.New("missing token")
	}

	url := API + "/bot" + n.Token + "/sendMessage"

	payload := new(bytes.Buffer)
	if err := json.NewEncoder(payload).Encode(n); err != nil {
		return err
	}

	resp, err := n.Client.Post(url, "application/json", payload)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var r apiResponse

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	if !r.OK {
		return errors.New(r.Description)
	}

	return nil
}
