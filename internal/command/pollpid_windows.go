package command

import (
	"errors"
	"time"
)

func pollPID(_ int, _ time.Duration) error {
	return errors.New("pwatch not supported on this platform")
}
