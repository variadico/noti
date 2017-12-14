// +build !darwin
// +build !windows

package command

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/variadico/noti/service/espeak"
	"github.com/variadico/noti/service/freedesktop"
)

func getBanner(title, message string, v *viper.Viper) notification {
	return &freedesktop.Notification{
		Summary:       title,
		Body:          message,
		ExpireTimeout: 500,
		AppIcon:       "utilities-terminal",
	}
}

func getSpeech(title, message string, v *viper.Viper) notification {
	return &espeak.Notification{
		Text:      fmt.Sprintf("%s %s", title, message),
		VoiceName: v.GetString("espeak.voiceName"),
	}
}
