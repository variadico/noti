package noti

import (
	"fmt"
	"os/exec"
)

type Notification struct {
	Summary  string
	Body     string
	IconName string
}

func (n *Notification) GetTitle() string {
	return n.Summary
}

func (n *Notification) SetTitle(t string) {
	n.Summary = t
}

func (n *Notification) GetMessage() string {
	return n.Body
}

func (n *Notification) SetMessage(m string) {
	n.Body = m
}

func (n *Notification) Notify() error {
	out, err := exec.Command("notify-send", n.Summary, n.Body).CombinedOutput()
	if err != nil {
		return fmt.Errorf("notify-send failed: %s\nOutput: %s", err, out)
	}

	return nil
}
