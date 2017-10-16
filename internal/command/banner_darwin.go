package command

import (
	"github.com/variadico/noti/service"
	"github.com/variadico/noti/service/nsuser"
)

func getBanner(title, message, sound string) service.Notification {
	return &nsuser.Notification{
		Title:           title,
		InformativeText: message,
		SoundName:       sound,
	}
}
