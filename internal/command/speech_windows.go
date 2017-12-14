package command

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/variadico/noti/service/speechsynthesizer"
)

func setSpeechDefaults(v *viper.Viper) {
	defaults := map[string]string{
		"speechsynthesizer.voice": "Microsoft David Desktop",
	}
	for key, val := range defaults {
		v.SetDefault(key, val)
	}

	envs := map[string]string{
		"speechsynthesizer.voice": "NOTI_VOICE",
	}
	for key, val := range envs {
		v.BindEnv(key, val)
	}
}

func getSpeech(title, message string, v *viper.Viper) notification {
	voice := v.GetString("speechsynthesizer.voice")

	return &speechsynthesizer.Notification{
		Text:  fmt.Sprintf("%s %s", title, message),
		Rate:  3,
		Voice: voice,
	}
}
