// +build windows

package speech

import (
	"github.com/variadico/noti"
)

// Notify via speech output is not supported on Windows yet, but we need it to be able to build the binary
func Notify(n noti.Params) error {
	return nil
}
