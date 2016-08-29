package main

import (
	"errors"
	"time"
)

func pollPID(pid int, interval time.Duration) error {
	return errors.New("pwatch not supported on this platform")
}
