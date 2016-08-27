// +build windows

package main

import (
	"time"
)

// Not supported on Windows yet.
func pollPID(pid int, interval time.Duration) error {
	return nil
}
