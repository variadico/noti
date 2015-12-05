// +build linux freebsd

// Package espeak speaks a notification using the espeak command on Linux and
// FreeBSD.
package espeak

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
		cmd = exec.Command("espeak", "--", n.Message)
	} else {
		cmd = exec.Command("espeak", "-v", n.Voice, "--", n.Message)
	}

	return cmd.Run()
}
