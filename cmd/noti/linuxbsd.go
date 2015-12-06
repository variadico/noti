// +build linux freebsd

package main

import (
	"flag"
	"fmt"

	"github.com/variadico/noti"
	"github.com/variadico/noti/espeak"
	"github.com/variadico/noti/libnotify"
)

const usageLinuxBSD = `
    -i, -icon
        Set icon name. You can pass a name from /usr/share/icons/gnome/32x32/ or
        /usr/share/notify-osd/icons/. Alternatively, you can specify a full
        filepath.
    -V, -voice
        Set voice.`

var (
	icon  *string
	voice *string
)

func init() {
	icon = flag.String("i", "", "")
	flag.StringVar(icon, "icon", "", "")
	voice = flag.String("V", "", "")
	flag.StringVar(voice, "voice", "", "")

	usageText = fmt.Sprintf(usageTmpl, usageLinuxBSD)
}

func notify(title, message string) error {
	var nt noti.NotifierMessager

	if *voice != "" {
		nt = &espeak.Notification{
			Voice:   *voice,
			Message: message,
		}

		message = fmt.Sprintf("%s %s", title, message)
	} else {
		nt = &libnotify.Notification{
			Summary:  title,
			Body:     message,
			IconName: *icon,
		}
	}

	nt.SetMessage(message)

	return nt.Notify()
}
