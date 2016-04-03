package noti

// Params are parameters that are passed to a notification's Notify function.
type Params struct {
	Title   string
	Message string

	Failure bool
	API     string
	Config  EnvGetter
}
