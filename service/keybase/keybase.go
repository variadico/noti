package keybase

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// KeybaseBin is cli binary
const KeybaseBin = "keybase"

// Errors when parsing notification settings
var (
	ErrorMissingConversation = errors.New("keybase: missing conversation (team or username) to send to")
	ErrorMissingMessage      = errors.New("keybase: missing message content")
	ErrorBadExplodingTime    = errors.New("keybase: explodingLifetime must be a time between 30s and 168h0m0s")
)

// Notification is a Keybase notification.
type Notification struct {
	// Conversation is the team name or users (comma-separated) to notify.
	Conversation string
	// ChannelName is the team's chat channel to send to. If empty, the team's
	// default channel will be used (typically "general").
	ChannelName string
	// Public toggles broadcasting a message to everyone (when Conversation is your
	// username), or to teams (when Conversation is your team name).
	Public bool
	// ExplodingLifetime will delete the message in the time given.
	ExplodingLifetime time.Duration
	// Message contents to be sent
	Message string
}

// prepareArgs builds the `keybase` cli arguments from the Notification settings
func prepareArgs(n *Notification) ([]string, error) {
	switch {
	case n.Conversation == "":
		return nil, ErrorMissingConversation
	case n.Message == "":
		return nil, ErrorMissingMessage
	case n.ExplodingLifetime < 0:
		return nil, ErrorBadExplodingTime
	}

	args := []string{"chat", "send"}
	if n.ChannelName != "" {
		args = append(args, "--channel", n.ChannelName)
	}
	if n.Public {
		args = append(args, "--public")
	}
	if n.ExplodingLifetime > 0 {
		args = append(args, "--exploding-lifetime", fmt.Sprint(n.ExplodingLifetime))
	}
	args = append(args, n.Conversation, n.Message)
	return args, nil
}

// Send triggers a Keybase notification.
func (n *Notification) Send() error {
	args, err := prepareArgs(n)
	if err != nil {
		return err
	}

	cmd := exec.Command(KeybaseBin, args...)
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("keybase: command error: %v", err)
	}
	return nil
}
