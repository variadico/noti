// +build linux freebsd

package main

import (
	"fmt"
	"os/exec"
)

// foreground and sound are only implemented on darwin
func notify(title, mesg, _ string, _, pbullet, cluster bool) error {
	if cluster == false {
		out, err := exec.Command("notify-send", title, mesg).CombinedOutput()
		if err != nil {
			return fmt.Errorf("notify-send failed: %s\nOutput: %s", err, out)
		}
	}
	if pbullet {
		return pbulletNotify(title, mesg)
	}
	return nil
}
