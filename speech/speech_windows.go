package speech

import (
	"errors"

	"github.com/variadico/noti"
)

// Notify speaks a notification. This will always return an error on Windows.
func Notify(n noti.Params) error {
	return errors.New("speech notification not supported on this platform")
}
