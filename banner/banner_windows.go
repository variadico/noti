package banner

import (
	"errors"

	"github.com/variadico/noti"
)

// Notify displays a notification. This will always return an error on Windows.
func Notify(n noti.Params) error {
	return errors.New("banner notification not supported on this platform")
}
