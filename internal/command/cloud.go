package command

import (
	"fmt"
	"net/http"
	"time"

	"github.com/variadico/noti/service"
	"github.com/variadico/noti/service/bearychat"
	"github.com/variadico/noti/service/hipchat"
	"github.com/variadico/noti/service/pushbullet"
	"github.com/variadico/noti/service/pushover"
	"github.com/variadico/noti/service/pushsafer"
	"github.com/variadico/noti/service/simplepush"
	"github.com/variadico/noti/service/slack"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

func getBearyChat(title, message, uri string) service.Notification {
	return &bearychat.Notification{
		Text:            fmt.Sprintf("**%s**\n%s", title, message),
		IncomingHookURI: uri,
		Client:          httpClient,
	}
}

func getHipChat(title, message, token, dest string) service.Notification {
	return &hipchat.Notification{
		Message:       fmt.Sprintf("%s\n%s", title, message),
		MessageFormat: "text",
		Client:        httpClient,
	}
}

func getPushbullet(title, message, token string) service.Notification {
	return &pushbullet.Notification{
		Title:  title,
		Body:   message,
		Type:   "note",
		Client: httpClient,
	}
}

func getPushover(title, message, token, user string) service.Notification {
	return &pushover.Notification{
		Title:   title,
		Message: message,
		Token:   token,
		User:    user,
		Client:  httpClient,
	}
}

func getPushsafer(title, message, key string) service.Notification {
	return &pushsafer.Notification{
		Title:      title,
		Message:    message,
		PrivateKey: key,
	}
}

func getSimplepush(title, message, key, event string) service.Notification {
	return &simplepush.Notification{
		Title:   title,
		Message: message,
		Key:     key,
		Event:   event,
	}
}

func getSlack(title, message, token, channel string) service.Notification {
	return &slack.Notification{
		Token:     token,
		Channel:   channel,
		Username:  "noti",
		Text:      fmt.Sprintf("%s\n%s", title, message),
		IconEmoji: ":rocket:",

		Client: httpClient,
	}
}
