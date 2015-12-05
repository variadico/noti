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
        Set notification sound. Default is Ping. Possible options are Basso,
        Blow, Bottle, Frog, Funk, Glass, Hero, Morse, Ping, Pop, Purr, Sosumi,
        Submarine, Tink. Check /System/Library/Sounds for available sounds.
    -V, -voice
        Set voice. Check System Preferences > Dictation & Speech for available
        voices.`

var (
	sound *string
	voice *string
)

func init() {
	sound = flag.String("s", "Ping", "")
	flag.StringVar(sound, "sound", "Ping", "")
	voice = flag.String("V", "", "")
	flag.StringVar(voice, "voice", "", "")

	usageText = fmt.Sprintf(usageTmpl, usageOSX)
}

func notify(title, message string) error {
	var nt noti.NotifierMessager

	if *voice != "" {
		nt = &nsspeechsynthesizer.Notification{
			Voice: *voice,
		}

		message = fmt.Sprintf("%s %s", title, message)
	} else {
		nt = &nsuser.Notification{
			Title:     title,
			SoundName: *sound,
		}
	}

	nt.SetMessage(message)

	return nt.Notify()
}
