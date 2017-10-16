package service

// Notification is the interface for all notifications.
type Notification interface {
	SetMessage(string)
	GetMessage() string
	Send() error
}
