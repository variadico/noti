//go:build windows
// +build windows

package command

import (
	"fmt"
	"syscall"
	"time"
)

func pollPID(pid int, _ time.Duration) error {
	kernel32 := syscall.NewLazyDLL("kernel32")
	openProcess := kernel32.NewProc("OpenProcess")
	waitForSingleObject := kernel32.NewProc("WaitForSingleObject")

	const (
		synchronize = 0x00100000
		infinite    = 0xFFFFFFFF
	)

	hProcess, _, lastErr := openProcess.Call(synchronize, 0, uintptr(pid))
	if hProcess == 0 {
		return fmt.Errorf("proc OpenProcess failed: %v", lastErr)
	}

	result, _, lastErr := waitForSingleObject.Call(hProcess, infinite)
	if result != 0 {
		return fmt.Errorf("proc WaitForSingleObject failed: %d %v", result, lastErr)
	}

	return nil
}
