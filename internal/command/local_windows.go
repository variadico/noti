package command

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/variadico/noti/service/notifyicon"
	"github.com/variadico/noti/service/speechsynthesizer"
)

func getBanner(title, message, _ string) notification {
	nt := &notifyicon.Notification{
		BalloonTipTitle: title,
		BalloonTipText:  message,
		BalloonTipIcon:  notifyicon.BalloonTipIconInfo,
	}

	return nt
}

func getSpeech(title, message string, v *viper.Viper) notification {
	voice := v.GetString("speechsynthesizer.voice")

	return &speechsynthesizer.Notification{
		Text:  fmt.Sprintf("%s %s", title, message),
		Rate:  3,
		Voice: voice,
	}
}
