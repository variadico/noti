package ntfy

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	n := Notification{
		URL:     "https://ntfy.sh/",
		Token:   "",
		Topic:   "topic",
		Message: "Message body test",
		Title:   "Test Message",
		Client:  &http.Client{Timeout: 3 * time.Second},
	}

	var mockResp apiResponse
	var hitServer bool

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		hitServer = true

		if r.Method != "POST" {
			t.Error("HTTP method should be POST")
		}

		var req Notification
		json.NewDecoder(r.Body).Decode(&req)

		if req.Topic == "" {
			t.Error("missing topic")
		}

		if req.Message == "" {
			t.Error("missing message")
		}

		if req.Title == "" {
			t.Error("missing title")
		}

		json.NewEncoder(rw).Encode(mockResp)
	}))
	defer ts.Close()

	n.URL = ts.URL

	if err := n.Send(); err != nil {
		t.Error(err)
	}

	if !hitServer {
		t.Error("didn't reach server")
	}
}
