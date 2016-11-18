package simplepush

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/variadico/noti"
)

const (
	keyEnv  = "NOTI_SIMPLEPUSH_KEY"
	eventEnv = "NOTI_SIMPLEPUSH_EVENT"

	// API is the Simplepush API endpoint.
	API = "https://api.simplepush.io/send"
)

var (
	errNoKey = noti.ConfigErrror{Env: keyEnv, Reason: "missing"}
)

type apiResponse struct {
	Status  string   `json:"status"`
	Errors  []string `json:"message"`
}

type configuration struct {
	key string
	event string
}

func envConfig(env noti.EnvGetter) (configuration, error) {
	key := env.Get(keyEnv)
	if key == "" {
		return configuration{}, errNoKey
	}

	event := env.Get(eventEnv)

	return configuration{
		key: key,
		event: event,
	}, nil
}

// Notify sends a push request to the Simplepush API.
func Notify(n noti.Params) error {
	config, err := envConfig(n.Config)
	if err != nil {
		return err
	}

	vals := make(url.Values)
	vals.Set("key", config.key)
	vals.Set("msg", n.Message)
	vals.Set("title", n.Title)
	if config.event != "" {
		vals.Set("event", config.event)
	}

	webClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := webClient.PostForm(n.API, vals)
	if err != nil {
		return err
	}
	defer resp.Body.Close()


	var r apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return fmt.Errorf("decoding response: %s", err)
	}

	if r.Status != "OK" {
		return noti.APIError{Site: "Simplepush", Msg: strings.Join(r.Errors, ": ")}
	}

	return nil
}
