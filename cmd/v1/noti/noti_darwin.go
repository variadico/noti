package main

import (
	"fmt"
	"os/exec"
)

const (
	activateReopen = `tell application "Terminal"
	activate
	reopen
end tell`

	displayNotification = "display notification %q with title %q sound name %q"
)

// notify displays a notification in OS X's notification center with a given
// title, message, and sound.
func notify(title, mesg, sound string, foreground, pbullet bool) error {
	if foreground {
		cmd := exec.Command("osascript", "-e", activateReopen)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	if pbullet {
		return pbulletNotify(title, mesg)
	}

	script := fmt.Sprintf(displayNotification, mesg, title, sound)
	cmd := exec.Command("osascript", "-e", script)
	return cmd.Run()
}
