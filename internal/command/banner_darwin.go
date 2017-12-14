package command

import (
	"github.com/spf13/viper"
	"github.com/variadico/noti/service/nsuser"
)

func setBannerDefaults(v *viper.Viper) {
	defaults := map[string]string{
		"nsuser.soundName":       "Ping",
		"nsuser.soundNameFail":   "Basso",
		"nsuser.informativeText": "Done!",
	}
	for key, val := range defaults {
		v.SetDefault(key, val)
	}

	envs := map[string]string{
		"nsuser.soundName":     "NOTI_SOUND",
		"nsuser.soundNameFail": "NOTI_SOUND_FAIL",
	}
	for key, val := range envs {
		v.BindEnv(key, val)
	}
}

func getBanner(title, message, sound string) notification {
	return &nsuser.Notification{
		Title:           title,
		InformativeText: message,
		SoundName:       sound,
	}
}
