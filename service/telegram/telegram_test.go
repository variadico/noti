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
		Client:  &http.Client{Timeout: 3 * time.Second},
	}

	var mockResp apiResponse
	var hitServer bool

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		hitServer = true

		var n Notification
		if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
			t.Error(err)
		}

		if r.Method != "POST" {
			t.Error("HTTP method should be POST.")
		}

		if n.ChatID == "" {
			t.Error("bot is not a member of the ", n.ChatID)
		}

		if n.Message == "" {
			t.Error("missing message")
		}

		if err := json.NewEncoder(rw).Encode(mockResp); err != nil {
			t.Error(err)
		}
	}))
	defer ts.Close()

	API = ts.URL

	// successful
	mockResp.OK = true
	if err := n.Send(); err != nil {
		t.Error(err)
	}

	if !hitServer {
		t.Error("didn't reach server")
	}

	// failure
	mockResp.OK = false
	if err := n.Send(); err == nil {
		t.Error("unexpected success")
	}
}
