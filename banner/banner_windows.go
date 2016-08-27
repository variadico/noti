package banner

import (
	toast "github.com/jacobmarshall/go-toast"
	"github.com/variadico/noti"
)

// Notify displays a Windows 10 Toast Notification.
func Notify(n noti.Params) error {
	notification := toast.Notification{
		AppID:   "noti",
		Title:   n.Title,
		Message: n.Message,
		Icon:    "",
		Actions: nil}

	return notification.Push()
}
