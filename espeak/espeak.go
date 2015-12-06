// +build linux freebsd

// Package espeak speaks a notification using the espeak command on Linux and
// FreeBSD.
//
// In order to use this package, you'll need to have espeak installed.
// On Ubuntu, espeak can be installed this way.
//    sudo apt-get install espeak
//
// To view a list of available voices, use this command.
//    espeak --voices
package espeak

import "os/exec"

// Notification is an espeak notification.
type Notification struct {
	Message string
	Voice   string
}

// GetMessage returns a notification's message.
func (n *Notification) GetMessage() string {
	return n.Message
}

// SetMessage sets a notification's message.
func (n *Notification) SetMessage(m string) {
	n.Message = m
}

// Notify speaks a notification's message. If the voice field is set, then it'll
// use that voice. It'll return an error if there was a problem executing
// espeak.
func (n *Notification) Notify() error {
	var cmd *exec.Cmd
	if n.Voice == "" {
		cmd = exec.Command("espeak", "--", n.Message)
	} else {
		cmd = exec.Command("espeak", "-v", n.Voice, "--", n.Message)
	}

	return cmd.Run()
}
