// +build darwin

package main

import (
	"flag"
	"fmt"

	"github.com/variadico/noti"
	"github.com/variadico/noti/nsspeechsynthesizer"
	"github.com/variadico/noti/nsuser"
)

const usageOSX = `
    -s, -sound
        Set notification sound.`

var (
	sound  *string
	speech *bool
)

func init() {
	sound = flag.String("s", "Ping", "")
	flag.StringVar(sound, "sound", "Ping", "")
	speech = flag.Bool("S", false, "")
	flag.BoolVar(speech, "speech", false, "")

	usageText = fmt.Sprintf(usageTmpl, usageOSX)
}

func notify(title, message string) error {
	var nt noti.Notifier

	if *speech {
		nt = &nsspeechsynthesizer.Notification{}
	} else {
		nt = &nsuser.Notification{}
	}

	nt.SetTitle(title)
	nt.SetMessage(message)

	return nt.Notify()
}
