package command

import (
	"fmt"
	"net/http"
	"time"

	"github.com/variadico/noti/service/bearychat"
	"github.com/variadico/noti/service/hipchat"
	"github.com/variadico/noti/service/pushbullet"
	"github.com/variadico/noti/service/pushover"
	"github.com/variadico/noti/service/pushsafer"
	"github.com/variadico/noti/service/simplepush"
	"github.com/variadico/noti/service/slack"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

func getBearyChat(title, message, uri string) notification {
	return &bearychat.Notification{
		Text:            fmt.Sprintf("**%s**\n%s", title, message),
		IncomingHookURI: uri,
		Client:          httpClient,
	}
}

func getHipChat(title, message, token, dest string) notification {
	return &hipchat.Notification{
		Message:       fmt.Sprintf("%s\n%s", title, message),
		MessageFormat: "text",
		Client:        httpClient,
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

func getPushover(title, message, token, user string) notification {
	return &pushover.Notification{
		Title:   title,
		Message: message,
		Token:   token,
		User:    user,
		Client:  httpClient,
	}
}

func getPushsafer(title, message, key string) notification {
	return &pushsafer.Notification{
		Title:      title,
		Message:    message,
		PrivateKey: key,
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
