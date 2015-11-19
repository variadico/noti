// Package noti describes a basic notification.
package noti

// Notifier describes a basic notification and can trigger a notification. At
// minimum, a Notifier should have a notification title and message.
type Notifier interface {
	// GetTitle returns a notification's title.
	GetTitle() string
	// SetTitle sets a notification's title.
	SetTitle(string)
	// GetMessage returns a notification's message.
	GetMessage() string
	// SetMessage sets a notification's message.
	SetMessage(string)
	// Notify triggers a notification.
	Notify() error
}
