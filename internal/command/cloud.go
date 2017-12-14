package command

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"github.com/variadico/noti/service/bearychat"
	"github.com/variadico/noti/service/hipchat"
	"github.com/variadico/noti/service/pushbullet"
	"github.com/variadico/noti/service/pushover"
	"github.com/variadico/noti/service/pushsafer"
	"github.com/variadico/noti/service/simplepush"
	"github.com/variadico/noti/service/slack"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

func setBearyChatDefaults(v *viper.Viper) {
	defaults := map[string]string{
		"bearychat.incomingHookURI": "",
	}
	for key, val := range defaults {
		v.SetDefault(key, val)
	}

	envs := map[string]string{
		"bearychat.incomingHookURI": "NOTI_BC_INCOMING_URI",
	}
	for key, val := range envs {
		v.BindEnv(key, val)
	}
}

func getBearyChat(title, message, uri string) notification {
	return &bearychat.Notification{
		Text:            fmt.Sprintf("**%s**\n%s", title, message),
		IncomingHookURI: uri,
		Client:          httpClient,
	}
}

func setHipChatDefaults(v *viper.Viper) {
	defaults := map[string]string{
		"hipchat.token":       "",
		"hipchat.destination": "",
	}
	for key, val := range defaults {
		v.SetDefault(key, val)
	}

	envs := map[string]string{
		"hipchat.token":       "NOTI_BC_INCOMING_URI",
		"hipchat.destination": "NOTI_HIPCHAT_DEST",
	}
	for key, val := range envs {
		v.BindEnv(key, val)
	}
}

func getHipChat(title, message, token, dest string) notification {
	return &hipchat.Notification{
		Message:       fmt.Sprintf("%s\n%s", title, message),
		MessageFormat: "text",
		Client:        httpClient,
	}
}

func setPushbulletDefaults(v *viper.Viper) {
	defaults := map[string]string{
		"pushbullet.token": "",
	}
	for key, val := range defaults {
		v.SetDefault(key, val)
	}

	envs := map[string]string{
		"pushbullet.token": "NOTI_PUSHBULLET_TOK",
	}
	for key, val := range envs {
		v.BindEnv(key, val)
	}
}

func getPushbullet(title, message, token string) notification {
	return &pushbullet.Notification{
		Title:  title,
		Body:   message,
		Type:   "note",
		Token:  token,
		Client: httpClient,
	}
}

func setPushoverDefaults(v *viper.Viper) {
	defaults := map[string]string{
		"pushover.token": "",
		"pushover.user":  "",
	}
	for key, val := range defaults {
		v.SetDefault(key, val)
	}

	envs := map[string]string{
		"pushover.token": "NOTI_PUSHOVER_TOK",
		"pushover.user":  "NOTI_PUSHOVER_DEST",
	}
	for key, val := range envs {
		v.BindEnv(key, val)
	}
}

func getPushover(title, message, token, user string) notification {
	return &pushover.Notification{
		Title:   title,
		Message: message,
		Token:   token,
		User:    user,
		Client:  httpClient,
	}
}

func setPushsaferDefaults(v *viper.Viper) {
	defaults := map[string]string{
		"pushsafer.token": "",
	}
	for key, val := range defaults {
		v.SetDefault(key, val)
	}

	envs := map[string]string{
		"pushsafer.token": "NOTI_PUSHSAFER_KEY",
	}
	for key, val := range envs {
		v.BindEnv(key, val)
	}
}

func getPushsafer(title, message, key string) notification {
	return &pushsafer.Notification{
		Title:      title,
		Message:    message,
		PrivateKey: key,
	}
}

func setSimplepushDefaults(v *viper.Viper) {
	defaults := map[string]string{
		"simplepush.key":   "",
		"simplepush.event": "",
	}
	for key, val := range defaults {
		v.SetDefault(key, val)
	}

	envs := map[string]string{
		"simplepush.key":   "NOTI_SIMPLEPUSH_KEY",
		"simplepush.event": "NOTI_SIMPLEPUSH_EVENT",
	}
	for key, val := range envs {
		v.BindEnv(key, val)
	}
}

func getSimplepush(title, message, key, event string) notification {
	return &simplepush.Notification{
		Title:   title,
		Message: message,
		Key:     key,
		Event:   event,
	}
}

func setSlackDefaults(v *viper.Viper) {
	defaults := map[string]string{
		"slack.token":   "",
		"slack.channel": "",
	}
	for key, val := range defaults {
		v.SetDefault(key, val)
	}

	envs := map[string]string{
		"slack.token":   "NOTI_SLACK_TOK",
		"slack.channel": "NOTI_SLACK_DEST",
	}
	for key, val := range envs {
		v.BindEnv(key, val)
	}
}

func getSlack(title, message, token, channel string) notification {
	return &slack.Notification{
		Token:     token,
		Channel:   channel,
		Username:  "noti",
		Text:      fmt.Sprintf("%s\n%s", title, message),
		IconEmoji: ":rocket:",

		Client: httpClient,
	}
}
