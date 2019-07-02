package telegram

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTelegramSend(t *testing.T) {
	n := Notification{
		ChatID:  "notifire",
		Message: "Testing notification",
		Client:  &http.Client{Timeout: 3 * time.Second,},
	}

	var mockResp apiResponse
	var hitServer bool

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		hitServer = true

		if r.Method != "POST" {
			t.Error("HTTP method should be POST.")
		}

		if r.FormValue("chat_id") == "" {
			t.Error("bot is not a member of the ", n.ChatID)
		}
		if r.FormValue("text") != n.Message {
			t.Error("missing message")
		}

		_ = json.NewEncoder(rw).Encode(mockResp)
	}))

	defer ts.Close()

	mockResp.OK = true // successful
	if err := n.Send(); err != nil {
		t.Error(err)
	}

	if !hitServer {
		t.Error("didn't reach server")
	}

	mockResp.OK = false // failure
	if err := n.Send(); err == nil {
		t.Error("unexpected success")
	}
}
