// +build !darwin
// +build !windows

package command

import (
	"github.com/spf13/viper"
	"github.com/variadico/noti/service/freedesktop"
)

func setBannerDefaults(v *viper.Viper) {
	// No banner defaults.
}

func getBanner(title, message, _ string) notification {
	return &freedesktop.Notification{
		Summary:       title,
		Body:          message,
		ExpireTimeout: 500,
		AppIcon:       "utilities-terminal",
	}
}
