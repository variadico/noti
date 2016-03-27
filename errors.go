package noti

import "fmt"

type ConfigErrror struct {
	Env    string
	Reason string
}

func (e ConfigErrror) Error() string {
	return fmt.Sprintf("invalid configuration for %s: %s", e.Env, e.Reason)
}

type APIError struct {
	Site string
	Msg  string
}

func (e APIError) Error() string {
	return fmt.Sprintf("%s API: %s", e.Site, e.Msg)
}
