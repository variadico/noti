// +build windows

package banner

import (
	"github.com/variadico/noti"
	toast "github.com/jacobmarshall/go-toast"
    "log"
)

// Notify displays a Windows 10 Toast Notification
func Notify(n noti.Params) error {
	 notification := toast.Notification{
        AppID: "noti",
        Title: n.Title,
        Message: n.Message,
        Icon: "",
        Actions: nil}

    err := notification.Push()

    if err != nil {
        log.Fatalln(err)
    }
    
	return nil
}
