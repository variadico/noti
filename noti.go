package noti

type Notifier interface {
	// Notify triggers a notification.
	Notify() error
}

type Titler interface {
	// GetTitle returns a notification's title.
	GetTitle() string
	// SetTitle sets a notification's title.
	SetTitle(string)
}

type Messager interface {
	// GetMessage returns a notification's message.
	GetMessage() string
	// SetMessage sets a notification's message.
	SetMessage(string)
}
