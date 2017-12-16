package command

import (
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Configuration Precedence
// * viper.Set
// * flag
// * env
// * file
// * defaults

var baseDefaults = map[string]interface{}{
	"defaults": []string{"banner"},
	"message":  "Done!",

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

func setNotiDefaults(v *viper.Viper) {
	for key, val := range baseDefaults {
		v.SetDefault(key, val)
	}
}

var keyEnvBindings = map[string]string{
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

func bindNotiEnv(v *viper.Viper) {
	for key, val := range keyEnvBindings {
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

// configureApp merges together different configuration sources.
func configureApp(v *viper.Viper, flags *pflag.FlagSet) error {
	setNotiDefaults(v)
	bindNotiEnv(v)

	if err := setupConfigFile(v); err != nil {
		return err
	}

	if flags != nil {
		v.BindPFlag("message", flags.Lookup("message"))
	}

	return nil
}

func enabledFromSlice(defaults []string) map[string]bool {
	// defaults should come from viper, which should  have processed baseDefaults
	// and config file values.

	services := map[string]bool{
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

	for _, name := range defaults {
		// Check if name is in services to avoid bad names from getting added
		// to map.
		if _, ok := services[name]; ok {
			services[name] = true
		}
	}

	return services
}

func enabledFromFlags(flags *pflag.FlagSet) map[string]bool {
	services := map[string]bool{
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

	// Visit flags that have been set.
	flags.Visit(func(f *pflag.Flag) {
		// pflag normalizes false, f, 0 to "false".
		if f.Value.Type() == "bool" && f.Value.String() == "false" {
			// Skip bool flags that are set to false.
			return
		}

		// Ignore flags that aren't service names.
		if _, ok := services[f.Name]; ok {
			services[f.Name] = true
		}
	})

	return services
}

func enabledServices(v *viper.Viper, flags *pflag.FlagSet) map[string]struct{} {
	var services map[string]bool

	// Highest precedence.
	if n := flags.NFlag(); n != 0 {
		services = enabledFromFlags(flags)
	}

	if s := os.Getenv("NOTI_DEFAULT"); s != "" {
		services = enabledFromSlice(strings.Split(s, " "))
	}

	// Lowest precedence.
	if s := v.GetStringSlice("defaults"); len(s) != 0 {
		services = enabledFromSlice(s)
	}

	filtered := make(map[string]struct{})
	for service, enabled := range services {
		if enabled {
			filtered[service] = struct{}{}
		}
	}

	return filtered
}

func getNotifications(v *viper.Viper, services map[string]struct{}) []notification {
	title := v.GetString("title")
	message := v.GetString("message")

	var notis []notification

	if _, ok := services["banner"]; ok {
		notis = append(notis, getBanner(title, message, v))
	}

	if _, ok := services["speech"]; ok {
		notis = append(notis, getSpeech(title, message, v))
	}

	if _, ok := services["bearychat"]; ok {
		notis = append(notis, getBearyChat(title, message, v))
	}

	if _, ok := services["hipchat"]; ok {
		notis = append(notis, getHipChat(title, message, v))
	}

	if _, ok := services["pushbullet"]; ok {
		notis = append(notis, getPushbullet(title, message, v))
	}

	if _, ok := services["pushover"]; ok {
		notis = append(notis, getPushover(title, message, v))
	}

	if _, ok := services["pushsafer"]; ok {
		notis = append(notis, getPushsafer(title, message, v))
	}

	if _, ok := services["simplepush"]; ok {
		notis = append(notis, getSimplepush(title, message, v))
	}

	if _, ok := services["slack"]; ok {
		notis = append(notis, getSlack(title, message, v))
	}

	return notis
}
