package pushsafer

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

var (
	// API is the Pushsafer API endpoint.
	API = "https://www.pushsafer.com/api"
)

type apiResponse struct {
	Info    string   `json:"info"`
	Status  int      `json:"status"`
	Request string   `json:"request"`
	Errors  []string `json:"errors"`
	Token   string   `json:"token"`
}

// Notification is a Pushsafer notification.
type Notification struct {
	Title   string
	Message string
	// Key is a private key or an alias key.
	Key string

	Client *http.Client
}

// Send sends a Pushsafer notification.
func (n *Notification) Send() error {
	vals := make(url.Values)
	vals.Set("k", n.Key)
	vals.Set("m", n.Message)
	vals.Set("t", n.Title)

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
	}

	return nil
}
