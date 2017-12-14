// +build !darwin
// +build !windows

package command

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/variadico/noti/service/espeak"
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

func setSpeechDefaults(v *viper.Viper) {
	defaults := map[string]string{
		"espeak.voiceName": "english-us",
	}
	for key, val := range defaults {
		v.SetDefault(key, val)
	}

	envs := map[string]string{
		"espeak.voiceName": "NOTI_VOICE",
	}
	for key, val := range envs {
		v.BindEnv(key, val)
	}
}

func getSpeech(title, message string, v *viper.Viper) notification {
	voice := v.GetString("espeak.voiceName")

	nt := &espeak.Notification{
		Text:      fmt.Sprintf("%s %s", title, message),
		VoiceName: "english-us",
	}

	if voice != "" {
		nt.VoiceName = voice
	}

	return nt
}
