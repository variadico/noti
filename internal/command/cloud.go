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

func getBearyChat(title, message string, v *viper.Viper) notification {
	return &bearychat.Notification{
		Text:            fmt.Sprintf("**%s**\n%s", title, message),
		IncomingHookURI: v.GetString("bearychat.incomingHookURI"),
		Client:          httpClient,
	}
}

func getHipChat(title, message string, v *viper.Viper) notification {
	return &hipchat.Notification{
		AccessToken:   v.GetString("hipchat.accessToken"),
		Room:          v.GetString("hipchat.room"),
		Message:       fmt.Sprintf("%s\n%s", title, message),
		MessageFormat: "text",
		Client:        httpClient,
	}
}

func getPushbullet(title, message string, v *viper.Viper) notification {
	return &pushbullet.Notification{
		Title:       title,
		Body:        message,
		Type:        "note",
		AccessToken: v.GetString("pushbullet.accessToken"),
		Client:      httpClient,
	}
}

func getPushover(title, message string, v *viper.Viper) notification {
	return &pushover.Notification{
		Title:   title,
		Message: message,
		Token:   v.GetString("pushover.token"),
		User:    v.GetString("pushover.user"),
		Client:  httpClient,
	}
}

func getPushsafer(title, message string, v *viper.Viper) notification {
	return &pushsafer.Notification{
		Title:      title,
		Message:    message,
		PrivateKey: v.GetString("pushsafer.privateKey"),
	}
}

func getSimplepush(title, message string, v *viper.Viper) notification {
	return &simplepush.Notification{
		Title:   title,
		Message: message,
		Key:     v.GetString("simplepush.key"),
		Event:   v.GetString("simplepush.event"),
	}
}

func getSlack(title, message string, v *viper.Viper) notification {
	return &slack.Notification{
		Token:     v.GetString("slack.token"),
		Channel:   v.GetString("slack.channel"),
		Username:  v.GetString("slack.username"),
		Text:      fmt.Sprintf("%s\n%s", title, message),
		IconEmoji: ":rocket:",

		Client: httpClient,
	}
}
