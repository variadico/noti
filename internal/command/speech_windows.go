package command

import (
	"fmt"

	"github.com/variadico/noti/services/speechsynthesizer"
)

func getSpeech(title, message, voice string) service.Notification {
	return &speechsynthesizer.Notification{
		Text:  fmt.Sprintf("%s %s", title, message),
		Rate:  3,
		Voice: voice,
	}
}
