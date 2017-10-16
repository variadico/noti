package command

import "github.com/variadico/noti/service/notifyicon"

func getBanner(title, message, _ string) service.Notification {
	nt := &notifyicon.Notification{
		BalloonTipTitle: title,
		BalloonTipText:  message,
		BalloonTipIcon:  notifyicon.BalloonTipIconInfo,
	}

}
