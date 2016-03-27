package hipchat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/variadico/noti"
)

const (
	destEnv = "NOTI_HIPCHAT_DEST"
	tokEnv  = "NOTI_HIPCHAT_TOK"

	API = "https://api.hipchat.com/v2/room/%s/notification"
)

var (
	errNoTok  = noti.ConfigErrror{Env: tokEnv, Reason: "missing"}
	errNoDest = noti.ConfigErrror{Env: destEnv, Reason: "missing"}
)

type apiResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

type apiRequest struct {
	Message       string `json:"message"`
	MessageFormat string `json:"message_format"`
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

func Notify(n noti.Notification) error {
	config, err := envConfig(n.Config)
	if err != nil {
		return err
	}

	payload := new(bytes.Buffer)
	err = json.NewEncoder(payload).Encode(apiRequest{
		Message:       fmt.Sprintf("%s\n%s", n.Title, n.Message),
		MessageFormat: "text",
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf(n.API, config.destination), payload)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.accessToken))
	req.Header.Set("Content-Type", "application/json")

	webClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := webClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err == io.EOF {
		return nil
	} else if err != nil {
		return fmt.Errorf("decoding response: %s", err)
	}

	if m := r.Error.Message; m != "" {
		return noti.APIError{Site: "hipchat", Msg: m}
	}

	return nil
}
