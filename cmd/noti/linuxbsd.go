// +build linux freebsd

package main

import (
	"flag"
	"fmt"

	"github.com/variadico/noti/notify"
)

const usageLinuxBSD = `
    -i, -icon
        Set icon name.`

var (
	icon *string
)

func init() {
	icon = flag.String("i", "", "")
	flag.StringVar(icon, "icon", "", "")

	usageText = fmt.Sprintf(usageTmpl, usageLinuxBSD)
}

func notify(title, message string) error {
	nt := notify.Notification{
		Summary: title,
		Body:    message,
	}

	return nt.Notify()
}
