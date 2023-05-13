package gchat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
)

// Notification is a Google Chat notification.
type Notification struct {
	// AppURL is your Google Chat App's webhook URL
	AppURL string

	// Text is the notification's message.
	Message string

	// Title is the message title
	Title string

	// Template is the template that combines title & message in a single text
	Template string

	Client *http.Client
}

type gchatMessage struct {
	Text string `json:"text"`
}

// Send triggers a Google Chat notification.
func (n *Notification) Send() error {
	if n.AppURL == "" {
		return errors.New("GChat: missing App URL")
	}
	if n.Message == "" {
		return errors.New("GChat: missing message text")
	}

	if n.Template == "" {
		return errors.New("GChat: missing message template")
	}

	t, err := template.New("").Parse(n.Template)
	if err != nil {
		return fmt.Errorf("GChat: %w", err)
	}

	data := &bytes.Buffer{}
	err = t.Execute(data, map[string]string{
		"message": n.Message,
		"title":   n.Title,
	})

	if err != nil {
		return fmt.Errorf("GChat: template error %w", err)
	}

	gChatBuffer := &bytes.Buffer{}
	message := gchatMessage{Text: data.String()}
	if err := json.NewEncoder(gChatBuffer).Encode(message); err != nil {
		return fmt.Errorf("GChat: JSON encoding error %w", err)
	}

	resp, err := n.Client.Post(n.AppURL, "application/json", gChatBuffer)
	if err != nil {
		return fmt.Errorf("GChat: Post error %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("GChat: error reading response %w", err)
		}
		return fmt.Errorf("GChat: %d\n %s", resp.StatusCode, string(b))
	}

	return nil
}
