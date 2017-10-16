package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

func macOSSounds(pass, fail string) (string, string) {
	if pass == "" {
		pass = "Ping"
	}

	if fail == "" {
		fail = "Basso"
	}

	return pass, fail
}

// printEnv prints all of the environment variables used in noti.
func printEnv() {
	envs := []string{
		"NOTI_BC_INCOMING_URI",
		"NOTI_DEFAULT",
		"NOTI_HIPCHAT_DEST",
		"NOTI_HIPCHAT_TOK",
		"NOTI_PUSHBULLET_TOK",
		"NOTI_PUSHOVER_DEST",
		"NOTI_PUSHOVER_TOK",
		"NOTI_PUSHSAFER_KEY",
		"NOTI_SIMPLEPUSH_EVENT",
		"NOTI_SIMPLEPUSH_KEY",
		"NOTI_SLACK_DEST",
		"NOTI_SLACK_TOK",
		"NOTI_SOUND",
		"NOTI_SOUND_FAIL",
		"NOTI_VOICE",
	}

	for _, env := range envs {
		if val, set := os.LookupEnv(env); set {
			fmt.Printf("%s=%s\n", env, val)
		}
	}
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
