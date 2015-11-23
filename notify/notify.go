// +build linux freebsd

// package notify can be used to trigger NotifyNotifications on systems that
// have libnotify and a notification server installed.
package notify

/*
#include <libnotify/notify.h>

void Notify(const char* summary, const char* body, const char* icon) {
	notify_init("Hello world!");
	NotifyNotification* note = notify_notification_new(summary, body, icon);
	notify_notification_show(note, NULL);
	g_object_unref(G_OBJECT(note));
	notify_uninit();
}
*/
import "C"

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
	C.Notify(
		C.CString(n.Summary),
		C.CString(n.Body),
		C.CString(n.IconName),
	)

	return nil
}
