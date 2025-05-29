package ntfy

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type apiResponse struct {
	// Message identifier
	Id string `json:"id"`

	// Message date time as Unix time stamp
	Time uint64 `json:"time"`

	// Unix time stamp for when the message will be deleted
	Expires uint64 `json:"expires"`

	// Message type
	Event string `json:"event"`

	// Topic ID
	Topic string `json:"topic"`

	// Message title
	Title string `json:"title"`

	// Message body
	Message string `json:"message"`
}

type Notification struct {
	// Base Ntfy URL
	URL string

	// Ntfy bearer access token (https://docs.ntfy.sh/publish/#access-tokens)
	Token string

	// Ntfy topic to publish to
	Topic string `json:"topic"`

	// Message body
	Message string `json:"message"`

	// Message title
	Title string `json:"title"`

	Client *http.Client `json:"-"`
}

// Send sends a ntfy notification to the configured url + topic.
func (n *Notification) Send() error {
	if n.URL == "" {
		return errors.New("missing Ntfy url")
	}

	if n.Topic == "" {
		return errors.New("missing topic id")
	}

	payload, err := json.Marshal(n)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", n.URL, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if n.Token != "" {
		req.Header.Set("Authorization", "Bearer "+n.Token)
	}

	resp, err := n.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	return nil
}
