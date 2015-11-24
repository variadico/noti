// +build darwin

package main

import (
	"flag"
	"fmt"

	"github.com/variadico/noti/nsuser"
)

const usageOSX = `
    -s, -sound
        Set notification sound.`

var (
	sound *string
)

func init() {
	sound = flag.String("s", "Ping", "")
	flag.StringVar(sound, "sound", "Ping", "")

	usageText = fmt.Sprintf(usageTmpl, usageOSX)
}

func notify(title, message string) error {
	nt := nsuser.Notification{
		Title:           title,
		InformativeText: message,
	}

	return nt.Notify()
}
