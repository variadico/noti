// +build !windows

package main

import (
	"os"
	"syscall"
	"time"
)

func watchPID(pid int, d time.Duration) error {
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

		time.Sleep(d)
	}

	return nil
}
