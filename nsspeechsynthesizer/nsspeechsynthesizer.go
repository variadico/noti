// +build darwin

// Package nsspeechsynthesizer can be used to speak a notification through
// NSSpeechSynthesizer on OS X.
package nsspeechsynthesizer

import "os/exec"

type Notification struct {
	Message string
	Voice   string
}

func (n *Notification) GetMessage() string {
	return n.Message
}

func (n *Notification) SetMessage(m string) {
	n.Message = m
}

func (n *Notification) Notify() error {
	var cmd *exec.Cmd
	if n.Voice == "" {
		cmd = exec.Command("say", n.Message)
	} else {
		cmd = exec.Command("say", "-v", n.Voice, n.Message)
	}

	return cmd.Run()
}
