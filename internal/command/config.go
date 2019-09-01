package command

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/variadico/vbs"
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

	"hipchat.accessToken": "",
	"hipchat.room":        "",

	"pushbullet.accessToken": "",
	"pushbullet.deviceIden":  "",

	"pushover.apiToken": "",
	"pushover.userKey":  "",

	"pushsafer.key": "",

	"simplepush.key":   "",
	"simplepush.event": "",

	"slack.token":    "",
	"slack.channel":  "",
	"slack.username": "noti",

	"mattermost.username":        "noti",
	"mattermost.channel":         "",
	"mattermost.incomingHookURI": "",
	"mattermost.iconurl":         "",

	"telegram.chatId": "",
	"telegram.token": "",
}

func setNotiDefaults(v *viper.Viper) {
	for key, val := range baseDefaults {
		v.SetDefault(key, val)
	}
}

var keyEnvBindings = map[string]string{
	"nsuser.soundName":     "NOTI_NSUSER_SOUNDNAME",
	"nsuser.soundNameFail": "NOTI_NSUSER_SOUNDNAMEFAIL",

	"say.voice": "NOTI_SAY_VOICE",

	"espeak.voiceName": "NOTI_ESPEAK_VOICENAME",

	"speechsynthesizer.voice": "NOTI_SPEECHSYNTHESIZER_VOICE",

	"bearychat.incomingHookURI": "NOTI_BEARYCHAT_INCOMINGHOOKURI",

	"hipchat.accessToken": "NOTI_HIPCHAT_ACCESSTOKEN",
	"hipchat.room":        "NOTI_HIPCHAT_ROOM",

	"pushbullet.accessToken": "NOTI_PUSHBULLET_ACCESSTOKEN",
	"pushbullet.deviceIden":  "NOTI_PUSHBULLET_DEVICEIDEN",

	"pushover.apiToken": "NOTI_PUSHOVER_APITOKEN",
	"pushover.userKey":  "NOTI_PUSHOVER_USERKEY",

	"pushsafer.key": "NOTI_PUSHSAFER_KEY",

	"simplepush.key":   "NOTI_SIMPLEPUSH_KEY",
	"simplepush.event": "NOTI_SIMPLEPUSH_EVENT",

	"slack.appurl":   "NOTI_SLACK_APPURL",
	"slack.token":    "NOTI_SLACK_TOKEN",
	"slack.channel":  "NOTI_SLACK_CHANNEL",
	"slack.username": "NOTI_SLACK_USERNAME",

	"mattermost.username":        "NOTI_MATTERMOST_USERNAME",
	"mattermost.incomingHookURI": "NOTI_MATTERMOST_INCOMINGHOOKURI",
	"mattermost.channel":         "NOTI_MATTERMOST_CHANNEL",
	"mattermost.iconurl":         "NOTI_MATTERMOST_ICONURL",
	"mattermost.type":            "NOTI_MATTERMOST_TYPE",

	"telegram.chatId": "NOTI_TELEGRAM_CHATID",
	"telegram.token": "NOTI_TELEGRAM_TOKEN",
}

var keyEnvBindingsDeprecated = map[string]string{
	"NOTI_NSUSER_SOUNDNAME":          "NOTI_SOUND",
	"NOTI_NSUSER_SOUNDNAMEFAIL":      "NOTI_SOUND_FAIL",
	"NOTI_SAY_VOICE":                 "NOTI_VOICE",
	"NOTI_ESPEAK_VOICENAME":          "NOTI_VOICE",
	"NOTI_SPEECHSYNTHESIZER_VOICE":   "NOTI_VOICE",
	"NOTI_BEARYCHAT_INCOMINGHOOKURI": "NOTI_BC_INCOMING_URI",
	"NOTI_HIPCHAT_ACCESSTOKEN":       "NOTI_HIPCHAT_TOK",
	"NOTI_HIPCHAT_ROOM":              "NOTI_HIPCHAT_DEST",
	"NOTI_PUSHBULLET_ACCESSTOKEN":    "NOTI_PUSHBULLET_TOK",
	"NOTI_PUSHOVER_APITOKEN":         "NOTI_PUSHOVER_TOK",
	"NOTI_PUSHOVER_USERKEY":          "NOTI_PUSHOVER_DEST",
	"NOTI_SLACK_TOKEN":               "NOTI_SLACK_TOK",
	"NOTI_SLACK_CHANNEL":             "NOTI_SLACK_DEST",
}

func bindNotiEnv(v *viper.Viper) error {
	for key, val := range keyEnvBindings {
		if err := v.BindEnv(key, val); err != nil {
			return err
		}
	}

	// Map old deprecated env vars to new ones.
	for newEnv, oldEnv := range keyEnvBindingsDeprecated {
		v := os.Getenv(oldEnv)
		if v == "" {
			continue
		}

		fmt.Fprintf(os.Stderr, "Warning: %s is deprecated, use %s instead\n",
			oldEnv, newEnv)
		fmt.Fprintf(os.Stderr, "Remapping %s=%s to %s\n", oldEnv, v, newEnv)

		if err := os.Setenv(newEnv, v); err != nil {
			return err
		}
	}

	return nil
}

func setupConfigFile(fileFlag string, v *viper.Viper) error {
	viper.SupportedExts = []string{"yaml"}
	var configPaths []string

	if fileFlag != "" {
		configPaths = append(configPaths, fileFlag)
	}

	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfig == "" {
		xdgConfig = filepath.Join(os.ExpandEnv("$HOME"), ".config", "noti", "noti.yaml")
	} else {
		xdgConfig = filepath.Join(xdgConfig, "noti", "noti.yaml")
	}

	configPaths = append(configPaths,
		filepath.Join(".", ".noti.yaml"),
		xdgConfig,
	)

	var config io.Reader
	var errMsg []string
	for _, p := range configPaths {
		data, err := ioutil.ReadFile(p)
		if err != nil {
			errMsg = append(errMsg, err.Error())
			continue
		}

		config = bytes.NewReader(data)
		break
	}
	if config == nil {
		return fmt.Errorf("failed to read config file: %v", errMsg)
	}

	v.SetConfigType("yaml")
	return v.ReadConfig(config)
}

// configureApp merges together different configuration sources.
func configureApp(v *viper.Viper, flags *pflag.FlagSet) error {
	setNotiDefaults(v)

	if err := bindNotiEnv(v); err != nil {
		return err
	}

	// Don't care about this error, fileFlag can be blank.
	fileFlag, _ := flags.GetString("file")
	if err := setupConfigFile(fileFlag, v); err != nil {
		// Not the end of the world if we can't read the config file.
		vbs.Println(err)
	}

	if flags == nil {
		return nil
	}

	return v.BindPFlag("message", flags.Lookup("message"))
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
		"mattermost": false,
		"telegram":   false,
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

func hasServiceFlags(flags *pflag.FlagSet) bool {
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
		"mattermost": false,
		"telegram":   false,
	}

	flags.Visit(func(f *pflag.Flag) {
		if _, ok := services[f.Name]; ok {
			services[f.Name] = true
		}
	})

	for _, enabled := range services {
		if enabled {
			return true
		}
	}
	return false
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
		"mattermost": false,
		"telegram":   false,
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

	if hasServiceFlags(flags) {
		// Highest precedence.
		services = enabledFromFlags(flags)
	} else if s := os.Getenv("NOTI_DEFAULT"); s != "" {
		services = enabledFromSlice(strings.Split(s, " "))
	} else if s := v.GetStringSlice("defaults"); len(s) != 0 {
		// Lowest precedence.
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

	if _, ok := services["mattermost"]; ok {
		notis = append(notis, getMattermost(title, message, v))
	}

	if _, ok := services["telegram"]; ok {
		notis = append(notis, getTelegram(title, message, v))
	}

	return notis
}
