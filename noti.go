// Package noti defines common methods used by user-facing notifications.
package noti

// Notifier represents a notification that can be triggered.
type Notifier interface {
	Notify() error
}

// Messager represents a notification that can get and set a text-based message.
type Messager interface {
	GetMessage() string
	SetMessage(string)
}

// NotifierMessager groups Notify, GetMessage, and SetMessage.
type NotifierMessager interface {
	Notifier
	Messager
}

// Titler represents a notification that can get and set a title.
type Titler interface {
	GetTitle() string
	SetTitle(string)
}
