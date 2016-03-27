package noti

type Notification struct {
	Title   string
	Message string
	Failure bool
	API     string
	Config  EnvGetter
}
