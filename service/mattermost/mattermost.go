package mattermost

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ErrInvalidResponse is returned when the response is not the expected result
var ErrInvalidResponse = errors.New("mattermost: invalid error response")

// will be thrown if an error occurs
// apiErrorResponse defines all fields which will be send by mattermost
type apiErrorResponse struct {
	ID         string `json:"id"`
	Message    string `json:"message"`
	Details    string `json:"detailed_error"`
	StatusCode int    `json:"status_code"`
	RequestID  string `json:"request_id"`
}

// String returns a string from all fields
func (ar *apiErrorResponse) String() string {
	bts, _ := json.Marshal(ar)
	return string(bts)
}

// Notification is a Mattermost notification.
type Notification struct {
	Text     string `json:"text"`
	Channel  string `json:"channel,omitempty"`
	Username string `json:"username,omitempty"`
	IconURL  string `json:"icon_url,omitempty"`
	Type     string `json:"type,omitempty"`

	IncomingHookURI string       `json:"-"`
	Client          *http.Client `json:"-"`
}

// Send triggers a Mattermost notification.
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

	// Check response
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		// Compare to mattermost success answer "ok"
		if !bytes.Equal(body, []byte("ok")) {
			return ErrInvalidResponse
		}

		return nil
	}

	var errResp apiErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		return fmt.Errorf("mattermost decoding response: %s", err.Error())
	}

	return fmt.Errorf("mattermost response: %s", errResp.String())
}
