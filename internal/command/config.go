package command

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func setNotiDefaults(v *viper.Viper) {
	defaults := map[string]string{
		"nsuser.soundName":     "Ping",
		"nsuser.soundNameFail": "Basso",

		"say.voice": "Alex",

		"espeak.voiceName": "english-us",

		"speechsynthesizer.voice": "Microsoft David Desktop",

		"bearychat.incomingHookURI": "",

		"hipchat.token":       "",
		"hipchat.destination": "",

		"pushbullet.token": "",

		"pushover.token": "",
		"pushover.user":  "",

		"pushsafer.privateKey": "",

		"simplepush.key":   "",
		"simplepush.event": "",

		"slack.token":   "",
		"slack.channel": "",
	}

	for key, val := range defaults {
		v.SetDefault(key, val)
	}
}

func bindNotiEnv(v *viper.Viper) {
	envs := map[string]string{
		"nsuser.soundName":     "NOTI_SOUND",
		"nsuser.soundNameFail": "NOTI_SOUND_FAIL",

		"say.voice": "NOTI_VOICE",

		"espeak.voiceName": "NOTI_VOICE",

		"speechsynthesizer.voice": "NOTI_VOICE",

		"bearychat.incomingHookURI": "NOTI_BC_INCOMING_URI",

		"hipchat.token":       "NOTI_BC_INCOMING_URI",
		"hipchat.destination": "NOTI_HIPCHAT_DEST",

		"pushbullet.token": "NOTI_PUSHBULLET_TOK",

		"pushover.token": "NOTI_PUSHOVER_TOK",
		"pushover.user":  "NOTI_PUSHOVER_DEST",

		"pushsafer.privateKey": "NOTI_PUSHSAFER_KEY",

		"simplepush.key":   "NOTI_SIMPLEPUSH_KEY",
		"simplepush.event": "NOTI_SIMPLEPUSH_EVENT",

		"slack.token":   "NOTI_SLACK_TOK",
		"slack.channel": "NOTI_SLACK_DEST",
	}

	for key, val := range envs {
		v.BindEnv(key, val)
	}
}

func setupConfigFile(v *viper.Viper) error {
	viper.SupportedExts = []string{"yaml"}
	v.SetConfigName(".noti")

	v.AddConfigPath(".")
	v.AddConfigPath("$HOME")

	return v.ReadInConfig()
}

// readEnv populates the initial config map.
func readEnv(env string) map[string]bool {
	// Initially everything is off.
	config := map[string]bool{
		"banner":     false,
		"bearychat":  false,
		"hipchat":    false,
		"pushbullet": false,
		"pushover":   false,
		"pushsafer":  false,
		"simplepush": false,
		"slack":      false,
		"speech":     false,
	}

	if env == "" {
		return config
	}

	envDefaults := strings.Split(env, " ")
	for _, name := range envDefaults {
		if _, found := config[name]; found {
			config[name] = true
		}
	}

	return config
}

// readFlags overrides anything set in the environment.
func readFlags(flags *pflag.FlagSet, config map[string]bool, defaultWasSet bool) error {
	for name := range config {
		if !flags.Changed(name) {
			// Flag was not set by user.

			if !defaultWasSet {
				continue
			}
		}

		val, err := flags.GetBool(name)
		if err != nil {
			return err
		}
		config[name] = val
	}

	return nil
}
