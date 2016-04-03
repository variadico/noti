package pushover

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
	tokEnv  = "NOTI_PUSHOVER_TOK"
	destEnv = "NOTI_PUSHOVER_DEST"

	// API is the Pushover API endpoint.
	API = "https://api.pushover.net/1/messages.json"
)

var (
	errNoTok  = noti.ConfigErrror{Env: tokEnv, Reason: "missing"}
	errNoDest = noti.ConfigErrror{Env: destEnv, Reason: "missing"}
)

type apiResponse struct {
	Info    string   `json:"info"`
	Status  int      `json:"status"`
	Request string   `json:"request"`
	Errors  []string `json:"errors"`
	Token   string   `json:"token"`
}

type configuration struct {
	accessToken string
	destination string
}

func envConfig(env noti.EnvGetter) (configuration, error) {
	tok := env.Get(tokEnv)
	if tok == "" {
		return configuration{}, errNoTok
	}

	dest := env.Get(destEnv)
	if dest == "" {
		return configuration{}, errNoDest
	}

	return configuration{
		accessToken: tok,
		destination: dest,
	}, nil
}

// Notify sends a push request to the Pushover API.
func Notify(n noti.Params) error {
	config, err := envConfig(n.Config)
	if err != nil {
		return err
	}

	vals := make(url.Values)
	vals.Set("token", config.accessToken)
	vals.Set("user", config.destination)
	vals.Set("message", n.Message)
	vals.Set("title", n.Title)

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

	if r.Status != 1 {
		return noti.APIError{Site: "Pushover", Msg: strings.Join(r.Errors, ": ")}
	} else if strings.Contains(r.Info, "no active devices") {
		return noti.APIError{Site: "Pushover", Msg: r.Info}
	}

	return nil
}
