// +build linux freebsd

package main

import (
	"flag"
	"fmt"

	"github.com/variadico/noti/libnotify"
)

const usageLinuxBSD = `
    -i, -icon
	Set icon name. You can pass a name from /usr/share/icons/gnome/32x32/
	or /usr/share/notify-osd/icons/. Alternatively, you can specify a full
	filepath.`

var (
	icon *string
)

func init() {
	icon = flag.String("i", "", "")
	flag.StringVar(icon, "icon", "", "")

	usageText = fmt.Sprintf(usageTmpl, usageLinuxBSD)
}

func notify(title, message string) error {
	nt := &libnotify.Notification{
		Summary:  title,
		Body:     message,
		IconName: *icon,
	}

	return nt.Notify()
}
