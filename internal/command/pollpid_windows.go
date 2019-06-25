// +build windows

package command

import (
	"errors"
	"time"
	"strconv"
	"syscall"
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32")

	openProcess = kernel32.NewProc("OpenProcess")
	waitForSingleObject = kernel32.NewProc("WaitForSingleObject")
)

const _SYNCHRONIZE = 0x00100000
const _INFINITE = 0xFFFFFFFF

func pollPID(pid int, _ time.Duration) error {
	hProcess, _, lastError := openProcess.Call(_SYNCHRONIZE, 0, uintptr(pid))

	if hProcess == 0 {
		return errors.New("OpenProcess failed with " + lastError.Error())
	}

	result, _, lastError := waitForSingleObject.Call(hProcess, _INFINITE)

	if result != 0 {
		return errors.New("WaitForSingleObject failed with " + strconv.Itoa(int(result)) + ", last error " + lastError.Error())
	}

	return nil
}
