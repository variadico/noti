package noti

import "os"

type EnvGetter interface {
	Get(string) string
}

type OSEnv struct{}

func (e OSEnv) Get(v string) string {
	return os.Getenv(v)
}

type MockEnv map[string]string

func (e MockEnv) Get(v string) string {
	return e[v]
}
