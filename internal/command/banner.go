// +build !darwin
// +build !windows

package command

import (
	"github.com/variadico/noti/service"
	"github.com/variadico/noti/service/freedesktop"
)

func getBanner(title, message, _ string) service.Notification {
	return &freedesktop.Notification{
		Summary:       title,
		Body:          message,
		ExpireTimeout: 500,
		AppIcon:       "utilities-terminal",
	}
}
