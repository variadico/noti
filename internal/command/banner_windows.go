package command

import (
	"github.com/spf13/viper"
	"github.com/variadico/noti/service"
	"github.com/variadico/noti/service/notifyicon"
)

func setBannerDefaults(v *viper.Viper) {
	// No banner defaults.
}

func getBanner(title, message, _ string) service.Notification {
	nt := &notifyicon.Notification{
		BalloonTipTitle: title,
		BalloonTipText:  message,
		BalloonTipIcon:  notifyicon.BalloonTipIconInfo,
	}

	return nt
}
