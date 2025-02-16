package command

import (
	"fmt"
	"html"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"github.com/variadico/noti/service/bark"
	"github.com/variadico/noti/service/bearychat"
	"github.com/variadico/noti/service/chanify"
	"github.com/variadico/noti/service/gchat"
	"github.com/variadico/noti/service/keybase"
	"github.com/variadico/noti/service/mattermost"
	"github.com/variadico/noti/service/ntfy"
	"github.com/variadico/noti/service/pushbullet"
	"github.com/variadico/noti/service/pushover"
	"github.com/variadico/noti/service/pushsafer"
	"github.com/variadico/noti/service/simplepush"
	"github.com/variadico/noti/service/slack"
	"github.com/variadico/noti/service/telegram"
	"github.com/variadico/noti/service/twilio"
	"github.com/variadico/noti/service/zulip"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

func getBearyChat(title, message string, v *viper.Viper) notification {
	return &bearychat.Notification{
		Text:            fmt.Sprintf("**%s**\n%s", title, message),
		IncomingHookURI: v.GetString("bearychat.incomingHookURI"),
		Client:          httpClient,
	}
}

func getKeybase(title, message string, v *viper.Viper) notification {
	var explodeTime time.Duration
	if v.GetString("keybase.explodingLifetime") != "" {
		// Error handling: if explodingLifetime is set to a unparseable duration,
		// viper will assign it to zero. Replace with -1, which will cause an early
		// error, to ensure the command does not send a regular message on accident.
		// Keybase's exploding messages have stricter security guarantees.
		explodeTime = v.GetDuration("keybase.explodingLifetime")
		if explodeTime == 0 {
			explodeTime = -1
		}
	}

	return &keybase.Notification{
		Conversation:      v.GetString("keybase.conversation"),
		ChannelName:       v.GetString("keybase.channel"),
		Public:            v.GetBool("keybase.public"),
		ExplodingLifetime: explodeTime,
		Message:           fmt.Sprintf("**%s**\n%s", title, message),
	}
}

func getPushbullet(title, message string, v *viper.Viper) notification {
	return &pushbullet.Notification{
		Title:       title,
		Body:        message,
		Type:        "note",
		AccessToken: v.GetString("pushbullet.accessToken"),
		DeviceIden:  v.GetString("pushbullet.deviceIden"),
		Client:      httpClient,
	}
}

func getPushover(title, message string, v *viper.Viper) notification {
	return &pushover.Notification{
		Title:    title,
		Message:  message,
		APIToken: v.GetString("pushover.apiToken"),
		UserKey:  v.GetString("pushover.userKey"),
		Client:   httpClient,
	}
}

func getPushsafer(title, message string, v *viper.Viper) notification {
	return &pushsafer.Notification{
		Title:   title,
		Message: message,
		Key:     v.GetString("pushsafer.key"),
		Client:  httpClient,
	}
}

func getSimplepush(title, message string, v *viper.Viper) notification {
	return &simplepush.Notification{
		Title:   title,
		Message: message,
		Key:     v.GetString("simplepush.key"),
		Event:   v.GetString("simplepush.event"),
		Client:  httpClient,
	}
}

func getSlack(title, message string, v *viper.Viper) notification {
	text := fmt.Sprintf("%s\n%s", title, message)
	if title == v.GetString("slack.username") {
		text = message
	}

	return &slack.Notification{
		Token:     v.GetString("slack.token"),
		Channel:   v.GetString("slack.channel"),
		Username:  v.GetString("slack.username"),
		AppURL:    v.GetString("slack.appurl"),
		Text:      text,
		IconEmoji: ":rocket:",

		Client: httpClient,
	}
}

func getGChat(title, message string, v *viper.Viper) notification {
	return &gchat.Notification{
		Message:  message,
		Title:    title,
		Template: v.GetString("gchat.template"),
		AppURL:   v.GetString("gchat.appurl"),
		Client:   httpClient,
	}
}

func getMattermost(title, message string, v *viper.Viper) notification {
	return &mattermost.Notification{
		IncomingHookURI: v.GetString("mattermost.incomingHookURI"),
		Channel:         v.GetString("mattermost.channel"),
		Username:        v.GetString("mattermost.username"),
		Text:            fmt.Sprintf("**%s %s**\n%s", title, ":rocket:", message),
		IconURL:         v.GetString("mattermost.iconurl"),
		Type:            v.GetString("mattermost.type"),

		Client: httpClient,
	}
}

func getTelegram(title, message string, v *viper.Viper) notification {
	return &telegram.Notification{
		ChatID:  v.GetString("telegram.chatId"),
		Token:   v.GetString("telegram.token"),
		Topic:   v.GetString("telegram.topic"),
		Message: fmt.Sprintf("<b>%s %s</b>\n%s", html.EscapeString(title), "🚀:", message),

		Client: httpClient,
	}
}

func getZulip(title, message string, v *viper.Viper) notification {
	return &zulip.Notification{
		BotAPIKey:       v.GetString("zulip.key"),
		BotEmailAddress: v.GetString("zulip.botAddress"),
		Endpoint:        v.GetString("zulip.URI"),
		Content:         fmt.Sprintf("%s:%s", title, message),
		Type:            v.GetString("zulip.type"),
		To:              v.GetString("zulip.to"),
		Client:          httpClient,
	}
}

func getTwilio(title, message string, v *viper.Viper) notification {
	return &twilio.Notification{
		Content:    fmt.Sprintf("%s:%s", title, message),
		NumberTo:   v.GetString("twilio.numberTo"),
		NumberFrom: v.GetString("twilio.numberFrom"),
		AccountSid: v.GetString("twilio.accountSid"),
		AuthToken:  v.GetString("twilio.authToken"),
	}
}

func getChanify(title, message string, v *viper.Viper) notification {
	return &chanify.Notification{
		ChannelURL:        v.GetString("chanify.channelURL"),
		Text:              message,
		Title:             title,
		Sound:             v.GetBool("chanify.sound"),
		Priority:          v.GetInt("chanify.priority"),
		InterruptionLevel: v.GetString("chanify.interruptionLevel"),
		Client:            httpClient,
	}
}

func getNtfy(title, message string, v *viper.Viper) notification {
	return &ntfy.Notification{
		URL:     v.GetString("ntfy.url"),
		Topic:   v.GetString("ntfy.topic"),
		Title:   title,
		Message: message,
		Client:  httpClient,
	}
}

func getBark(title, message string, v *viper.Viper) notification {
	return &bark.Notification{
		URL:       v.GetString("bark.apiurl"),
		DeviceKey: v.GetString("bark.key"),
		Title:     title,
		Body:      message,
		Client:    httpClient,
	}
}
