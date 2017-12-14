package command

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/variadico/noti/service/nsuser"
	"github.com/variadico/noti/service/say"
)

func getBanner(title, message, sound string) notification {
	return &nsuser.Notification{
		Title:           title,
		InformativeText: message,
		SoundName:       sound,
	}
}

func getSpeech(title, message string, v *viper.Viper) notification {
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
