package pushbullet

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
)

const (
	// AccessTokenEnv should contain a user's Pushbullet access token.
	AccessTokenEnv = "PUSHBULLET_ACCESS_TOKEN"

	apiEndpoint = "https://api.pushbullet.com/v2/pushes"
	reqJSON     = `{"body":%q,"title":%q,"type":"note"}`
)

var (
	// ErrNoAccessToken is returned if no access token was provided.
	ErrNoAccessToken = errors.New("Pushbullet access token is missing")
)

// Notification is a Pushbullet notification.
type Notification struct {
	AccessToken string
	Title       string
	Body        string
}

// GetTitle returns a notification's title.
func (n *Notification) GetTitle() string {
	return n.Title
}

// SetTitle sets a notification's title.
func (n *Notification) SetTitle(t string) {
	n.Title = t
}

// GetMessage returns a notification's message.
func (n *Notification) GetMessage() string {
	return n.Body
}

// SetMessage sets a notification's message.
func (n *Notification) SetMessage(m string) {
	n.Body = m
}

// Notify displays a notification on all the devices a user has registered with
// Pushbullet.
func (n *Notification) Notify() error {
	if n.AccessToken == "" {
		return ErrNoAccessToken
	}

	payload := bytes.NewBuffer([]byte(fmt.Sprintf(reqJSON, n.Body, n.Title)))

	req, err := http.NewRequest("POST", apiEndpoint, payload)
	if err != nil {
		return err
	}
	req.Header.Set("Access-Token", n.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	if _, err = http.DefaultClient.Do(req); err != nil {
		return err
	}

	return nil
}
