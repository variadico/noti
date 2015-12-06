// +build darwin

// Package nsspeechsynthesizer speaks a notification using the say command on OS
// X.
//
// In order to use this package, you'll need to have say installed. It should be
// installed by default on OS X.
//
// To view a list of available voices, use this command.
//    say -v ?
package nsspeechsynthesizer

import "os/exec"

// Notification is a NSSpeechSynthesizer notification.
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
// use that voice. It'll return an error if there was a problem executing say.
func (n *Notification) Notify() error {
	var cmd *exec.Cmd
	if n.Voice == "" {
		cmd = exec.Command("say", n.Message)
	} else {
		cmd = exec.Command("say", "-v", n.Voice, n.Message)
	}

	return cmd.Run()
}
