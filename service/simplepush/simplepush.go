package simplepush

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

var (
	// API is the Simplepush API endpoint.
	API = "https://api.simplepush.io/send"
)

type apiResponse struct {
	Status string   `json:"status"`
	Errors []string `json:"message"`
}

// Notification is a simplepush notification.
type Notification struct {
	Key     string
	Message string
	Title   string
	Event   string

	Client *http.Client
}

// Send sends a simplepush notification.
func (n *Notification) Send() error {
	if n.Key == "" {
		return errors.New("missing key")
	}

	vals := make(url.Values)
	vals.Set("key", n.Key)
	vals.Set("title", n.Title)
	vals.Set("msg", n.Message)
	vals.Set("event", n.Event)

	resp, err := n.Client.PostForm(API, vals)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	if r.Status != "OK" {
		return errors.New(strings.Join(r.Errors, ": "))
	}

	return nil

}
