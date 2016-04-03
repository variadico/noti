package noti

// Notification is a noti notification. It contains information that's used in a
// notification's Notify function.
type Notification struct {
	Title   string
	Message string

	Failure bool
	API     string
	Config  EnvGetter
}
