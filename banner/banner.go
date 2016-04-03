// +build linux freebsd

package banner

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/variadico/noti"
)

// Notify displays a libnotify notification.
func Notify(n noti.Params) error {
	_, err := exec.LookPath("notify-send")
	if err != nil {
		return errors.New("Install 'notify-send' and try again")
	}

	cmd := exec.Command("notify-send", n.Title, n.Message)
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("Banner: %s", err)
	}

	return nil
}
