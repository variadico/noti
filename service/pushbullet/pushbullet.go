package pushbullet

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

var (
	// API is the Pushbullet API endpoint.
	API = "https://api.pushbullet.com/v2/pushes"
)

type apiResponse struct {
	Active                  bool    `json:"active"`
	Iden                    string  `json:"iden"`
	Created                 float64 `json:"created"`
	Modified                float64 `json:"modified"`
	Type                    string  `json:"type"`
	Dismissed               bool    `json:"dismissed"`
	Direction               string  `json:"direction"`
	SenderIden              string  `json:"sender_iden"`
	SenderEmail             string  `json:"sender_email"`
	SenderEmailNormalized   string  `json:"sender_email_normalized"`
	SenderName              string  `json:"sender_name"`
	ReceiverIden            string  `json:"receiver_iden"`
	ReceiverEmail           string  `json:"receiver_email"`
	ReceiverEmailNormalized string  `json:"receiver_email_normalized"`
	Title                   string  `json:"title"`
	Body                    string  `json:"body"`
	Error                   struct {
		Code    string `json:"code"`
		Type    string `json:"type"`
		Message string `json:"message"`
		Cat     string `json:"cat"`
	} `json:"error"`
	ErrorCode string `json:"error_code"`
}

// Notification is a pushbullet notification.
type Notification struct {
	Body  string `json:"body"`
	Title string `json:"title"`
	Type  string `json:"type"`

	Token  string       `json:"-"`
	Client *http.Client `json:"-"`
}

// Send sends a Pushbullet notification.
func (n *Notification) Send() error {
	payload := new(bytes.Buffer)
	if err := json.NewEncoder(payload).Encode(n); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", API, payload)
	if err != nil {
		return err
	}
	req.Header.Set("Access-Token", n.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	if r.ErrorCode != "" {
		return errors.New(r.ErrorCode)
	}

	return nil
}
