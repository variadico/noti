package nsspeechsynthesizer

import "os/exec"

type Notification struct {
	Title   string
	Message string
	Voice   string
}

func (n *Notification) GetTitle() string {
	return n.Title
}

func (n *Notification) SetTitle(t string) {
	n.Title = t
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
		cmd = exec.Command("say", n.Title, n.Message)
	} else {
		cmd = exec.Command("say", "-v", n.Voice, n.Title, n.Message)
	}

	return cmd.Run()
}
