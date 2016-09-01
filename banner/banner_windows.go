// +build windows

package banner

import (
	toast "github.com/jacobmarshall/go-toast"
	"github.com/variadico/noti"
	"log"
)

// Notify displays a Windows 10 Toast Notification
func Notify(n noti.Params) {
	notification := toast.Notification{
		AppID:   "noti",
		Title:   n.Title,
		Message: n.Message,
		Icon:    "",
		Actions: nil}

	if err != nil {
		log.Fatalln(err)
	}

	return notification.Push()
}
