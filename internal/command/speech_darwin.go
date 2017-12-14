package command

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/variadico/noti/service"
	"github.com/variadico/noti/service/say"
)

func setSpeechDefaults(v *viper.Viper) {
	defaults := map[string]string{
		"say.voice": "Alex",
	}
	for key, val := range defaults {
		v.SetDefault(key, val)
	}

	envs := map[string]string{
		"say.voice": "NOTI_VOICE",
	}
	for key, val := range envs {
		v.BindEnv(key, val)
	}
}

func getSpeech(title, message string, v *viper.Viper) service.Notification {
	voice := v.GetString("say.voice")

	nt := &say.Notification{
		Voice: "Alex",
		Text:  fmt.Sprintf("%s %s", title, message),
		Rate:  200,
	}

	if voice != "" {
		nt.Voice = voice
	}

	return nt
}
