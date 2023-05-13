package chanify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Notification is a Chanify notification.
type Notification struct {
	// ChannelURL is the URL of the Chanify server
	// usually https://api.chanify.net/v1/sender/<TOKEN>
	ChannelURL string

	// Text is the notification's message.
	Text string

	// Title is the message title
	Title string

	// Sound is the notification's sound enabled if not in DnD
	// follow https://github.com/chanify/chanify#send-text
	Sound bool

	// Level of the notification used by iOS to display it
	// If unsure use 10 (that's the maximum)
	// follow https://github.com/chanify/chanify#send-text
	Priority int

	// Allows notification to bypass DnD or to be shown as TIME SENSITIVE
	// Use `active`, `passive`, `time-sensitive`
	// follow https://github.com/chanify/chanify#send-text
	InterruptionLevel string

	Client *http.Client
}

type chanifyMessage struct {
	Title             string `json:"title"`
	Text              string `json:"text"`
	Sound             bool   `json:"sound"`
	Priority          int    `json:"priority"`
	InterruptionLevel string `json:"interruption_level"`
}

// Send triggers a Chanify notification.
func (n *Notification) Send() error {
	if n.ChannelURL == "" {
		return errors.New("Chanify: missing channel URL")
	}
	if n.Text == "" {
		return errors.New("Chanify: missing message text")
	}

	message := chanifyMessage{
		Title:             n.Title,
		Text:              n.Text,
		Sound:             n.Sound,
		Priority:          n.Priority,
		InterruptionLevel: n.InterruptionLevel,
	}

	gChatBuffer := &bytes.Buffer{}
	if err := json.NewEncoder(gChatBuffer).Encode(message); err != nil {
		return fmt.Errorf("Chanify: JSON encoding error %w", err)
	}

	req, err := http.NewRequest("POST", n.ChannelURL, gChatBuffer)
	if err != nil {
		return fmt.Errorf("Chanify: error creating request %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.Client.Do(req)
	if err != nil {
		return fmt.Errorf("Chanify: Post error %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Chanify: error reading response %w", err)
		}
		return fmt.Errorf("Chanify: %d\n %s", resp.StatusCode, string(b))
	}

	return nil
}
