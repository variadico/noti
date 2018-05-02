package mattermost

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var ErrInvalidErrorResponse = errors.New("Invalid Error response")

type apiErrorResponse struct {
	Id         string `json:"id"`
	Message    string `json:"message"`
	Details    string `json:"detailed_error"`
	StatusCode int    `json:"status_code"`
	RequestId  string `json:"request_id"`
}

type Notification struct {
	Text     string `json:"text"`
	Channel  string `json:"channel,omitempty"`
	Username string `json:"username,omitempty"`
	IconURL  string `json:"icon_url,omitempty"`
	Type     string `json:"type,omitempty"`

	IncomingHookURI string       `json:"-"`
	Client          *http.Client `json:"-"`
}

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

	if resp.StatusCode != 200 {
		var errResp apiErrorResponse

		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return fmt.Errorf("decoding response: %s", err)
		}
		return fmt.Errorf("response: %s", errResp)
	}

	return nil
}
