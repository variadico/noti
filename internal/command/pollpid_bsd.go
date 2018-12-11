// +build freebsd netbsd openbsd dragonfly darwin

package command

import "unsafe"
import "strconv"
import "syscall"
import "time"
import "errors"
import "os"

import "golang.org/x/sys/unix"

func pollPID(pid int, _ time.Duration) error {
        kq, err := syscall.Kqueue()

        if err != nil {
                return err
        }

        if kq < 0 {
                return errors.New("kqueue failed with " + strconv.Itoa(-kq))
        }

        kqFile := os.NewFile(uintptr(kq), "kq")

        defer kqFile.Close()

        syscall.Syscall(syscall.SYS_FCNTL, uintptr(kq), syscall.F_SETFD, syscall.FD_CLOEXEC)

        ev := make([]syscall.Kevent_t, 1, 1)

        *(*int)(unsafe.Pointer(&ev[0].Ident)) = pid

        ev[0].Filter = unix.EVFILT_PROC
        ev[0].Flags = unix.EV_ADD
        ev[0].Fflags = unix.NOTE_EXIT
        ev[0].Data = 0
        ev[0].Udata = nil

        n, err := syscall.Kevent(kq, ev, nil, nil)

        if err != nil {
                return err
        }

        if n < 0 {
                return errors.New("kevent failed with " + strconv.Itoa(-n))
        }

        for {
                ev[0] = syscall.Kevent_t{}
                n, err = syscall.Kevent(kq, nil, ev, nil)

                if err != nil {
                        return err
                }

                if (ev[0].Fflags & unix.NOTE_EXIT) != 0 {
                        break
                }
        }

        return nil
}

