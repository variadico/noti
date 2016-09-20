package bearychat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/variadico/noti"
)

const (
	incomingEnv = "NOTI_BC_INCOMING_URI"
)

var (
	errNoIncoming = noti.ConfigErrror{Env: incomingEnv, Reason: "missing"}
)

type configuration struct {
	incoming string
}

type incomingPayload struct {
	Text string `json:"text"`
}

type incomingResponse struct {
	Code  int    `json:"code"`
	Error string `json:"string"`
}

func envConfig(env noti.EnvGetter) (configuration, error) {
	incoming := env.Get(incomingEnv)
	if incoming == "" {
		return configuration{}, errNoIncoming
	}

	return configuration{incoming}, nil
}

// Notify sends a message request to BearyChat's incoming hook.
func Notify(n noti.Params) error {
	config, err := envConfig(n.Config)
	if err != nil {
		return err
	}

	payload := new(bytes.Buffer)
	err = json.NewEncoder(payload).Encode(incomingPayload{
		Text: fmt.Sprintf("**%s**\n%s", n.Title, n.Message),
	})
	if err != nil {
		return err
	}

	webClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := webClient.Post(config.incoming, "application/json", payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r incomingResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return fmt.Errorf("decoding response: %s", err)
	}

	if r.Code != 0 {
		return noti.APIError{Site: "BearyChat", Msg: r.Error}
	}

	return nil
}
