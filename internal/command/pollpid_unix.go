// +build !windows

package command

import (
	"os"
	"syscall"
	"time"
)

func pollPID(pid int, interval time.Duration) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	for {
		// If I can send sig 0, then process is still alive.
		err = proc.Signal(syscall.Signal(0))
		if err != nil {
			break
		}

		time.Sleep(interval)
	}

	return nil
}
