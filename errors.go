package noti

import "fmt"

// ConfigErrror is a configuration error for noti packages.
type ConfigErrror struct {
	Env    string
	Reason string
}

func (e ConfigErrror) Error() string {
	return fmt.Sprintf("invalid configuration for %s: %s", e.Env, e.Reason)
}

// APIError is an API error that's returned if a notification API request
// failed.
type APIError struct {
	Site string
	Msg  string
}

func (e APIError) Error() string {
	return fmt.Sprintf("%s API: %s", e.Site, e.Msg)
}
