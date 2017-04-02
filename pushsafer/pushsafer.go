package pushsafer

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
	keyEnv  = "NOTI_PUSHSAFER_KEY"

	// API is the Pushsafer API endpoint.
	API = "https://www.pushsafer.com/api"
)

var (
	errNoKey  = noti.ConfigErrror{Env: keyEnv, Reason: "missing"}
)

type apiResponse struct {
	Info    string   `json:"info"`
	Status  int      `json:"status"`
	Request string   `json:"request"`
	Errors  []string `json:"errors"`
	Token   string   `json:"token"`
}

type configuration struct {
	privateKey string
}

func envConfig(env noti.EnvGetter) (configuration, error) {
	key := env.Get(keyEnv)
	if key == "" {
		return configuration{}, errNoKey
	}

	return configuration{
		privateKey: key,
	}, nil
}

// Notify sends a push request to the Pushsafer API.
func Notify(n noti.Params) error {
	config, err := envConfig(n.Config)
	if err != nil {
		return err
	}

	vals := make(url.Values)
	vals.Set("k", config.privateKey)
	vals.Set("m", n.Message)
	vals.Set("t", n.Title)

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
		return noti.APIError{Site: "Pushsafer", Msg: strings.Join(r.Errors, ": ")}
	} else if strings.Contains(r.Info, "no active devices") {
		return noti.APIError{Site: "Pushsafer", Msg: r.Info}
	}

	return nil
}
