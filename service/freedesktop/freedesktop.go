//go:build !darwin && !windows
// +build !darwin,!windows

package freedesktop

import (
	"fmt"

	dbus "github.com/godbus/dbus/v5"
)

// Notification is a Freedesktop notification.
type Notification struct {
	AppName    string
	ReplacesID uint
	AppIcon    string
	Summary    string
	Body       string
	Actions    []string
	// Hints         string
	ExpireTimeout int
}

// Send triggers a desktop notification.
func (n *Notification) Send() error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return fmt.Errorf("dbus connect: %s", err)
	}
	defer conn.Close()

	fdn := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")

	// 0 is a total magic number. ¯\_(ツ)_/¯
	resp := fdn.Call(
		"org.freedesktop.Notifications.Notify", 0,
		n.AppName,
		uint32(n.ReplacesID),
		n.AppIcon,
		n.Summary,
		n.Body,
		n.Actions,
		map[string]dbus.Variant{},
		int32(n.ExpireTimeout),
	)

	if resp.Err != nil {
		return fmt.Errorf("notify: %s", resp.Err)
	}

	return nil
}
