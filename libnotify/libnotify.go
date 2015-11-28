// +build linux freebsd

// Package libnotify can be used to display a NotifyNotification on systems
// that have libnotify and a notification server installed.
package libnotify

/*
#cgo pkg-config: --cflags glib-2.0
#cgo pkg-config: --cflags gtk+-2.0
#cgo LDFLAGS: -lnotify

#include <stdlib.h>
#include <errno.h>
#include <glib.h>
#include <libnotify/notify.h>

int notify(const char* summary, const char* body, const char* icon) {
	errno = 0;
	notify_init("noti");
	NotifyNotification* nt = notify_notification_new(summary, body, icon);

	notify_notification_set_timeout(nt, 3000);

	if (!notify_notification_show(nt, NULL)) {
		return 1;
	}

	g_object_unref(G_OBJECT(nt));
	notify_uninit();
	return 0;
}
*/
import "C"

import (
	"errors"
	"unsafe"
)

type Notification struct {
	Summary  string
	Body     string
	IconName string
}

func (n *Notification) GetTitle() string {
	return n.Summary
}

func (n *Notification) SetTitle(t string) {
	n.Summary = t
}

func (n *Notification) GetMessage() string {
	return n.Body
}

func (n *Notification) SetMessage(m string) {
	n.Body = m
}

func (n *Notification) Notify() error {
	s := C.CString(n.Summary)
	b := C.CString(n.Body)
	i := C.CString(n.IconName)
	defer C.free(unsafe.Pointer(s))
	defer C.free(unsafe.Pointer(b))
	defer C.free(unsafe.Pointer(i))

	rt, _ := C.notify(s, b, i)
	if rt == 1 {
		return errors.New("NotifyNotification show failed")
	}

	return nil
}
