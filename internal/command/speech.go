// +build !darwin
// +build !windows

package command

import (
	"fmt"

	"github.com/variadico/noti/service"
	"github.com/variadico/noti/service/espeak"
)

func getSpeech(title, message, voice string) service.Notification {
	nt := &espeak.Notification{
		Text:      fmt.Sprintf("%s %s", title, message),
		VoiceName: "english-us",
	}

	if voice != "" {
		nt.VoiceName = voice
	}

	return nt
}
