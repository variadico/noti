package banner

/*
// Compiler flags.
#cgo CFLAGS: -Wall -x objective-c -arch x86_64 -std=gnu99 -fobjc-arc
// Linker flags.
#cgo LDFLAGS: -framework Foundation -arch x86_64

#import "banner_darwin.h"
*/
import "C"
import (
	"unsafe"

	"github.com/variadico/noti"
)

const (
	soundEnv     = "NOTI_SOUND"
	soundFailEnv = "NOTI_SOUND_FAIL"
)

// Notify displays a NSUserNotification.
func Notify(n noti.Params) error {
	var sound string
	if n.Failure {
		sound = n.Config.Get(soundFailEnv)
		if sound == "" {
			sound = "Basso"
		}
	} else {
		sound = n.Config.Get(soundEnv)
		if sound == "" {
			sound = "Ping"
		}
	}

	t := C.CString(n.Title)
	m := C.CString(n.Message)
	s := C.CString(sound)
	defer C.free(unsafe.Pointer(t))
	defer C.free(unsafe.Pointer(m))
	defer C.free(unsafe.Pointer(s))

	C.BannerNotify(t, m, s)

	return nil
}
