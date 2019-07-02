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
		Token:   "token",
		Client:  &http.Client{Timeout: 3 * time.Second,},
	}

	var mockResp apiResponse
	var hitServer bool

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		hitServer = true

		decoder := json.NewDecoder(r.Body)

		var n Notification

		_ = decoder.Decode(&n)

		if r.Method != "POST" {
			t.Error("HTTP method should be POST.")
		}

		if n.ChatID == "" {
			t.Error("bot is not a member of the ", n.ChatID)
		}

		if n.Message == "" {
			t.Error("missing message")
		}

		_ = json.NewEncoder(rw).Encode(mockResp)
	}))

	defer ts.Close()

	API = ts.URL

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
