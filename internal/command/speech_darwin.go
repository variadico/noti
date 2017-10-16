package command

import (
	"fmt"

	"github.com/variadico/noti/service"
	"github.com/variadico/noti/service/say"
)

func getSpeech(title, message, voice string) service.Notification {
	nt := &say.Notification{
		Voice: "Alex",
		Text:  fmt.Sprintf("%s %s", title, message),
		Rate:  200,
	}

	if voice != "" {
		nt.Voice = voice
	}

	return nt
}
