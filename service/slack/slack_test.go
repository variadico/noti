package slack

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	n := Notification{
		Text:    "mesg",
		Token:   "tok",
		Channel: "chan",
		Client:  &http.Client{Timeout: 3 * time.Second},
	}
	var mockResp apiResponse
	var hitServer bool

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		hitServer = true

		if r.Method != "POST" {
			t.Error("HTTP method should be POST.")
		}

		if r.FormValue("token") == "" {
			t.Error("missing access token")
		}
		if r.FormValue("channel") == "" {
			t.Error("missing destination")
		}
		if r.FormValue("text") != n.Text {
			t.Error("missing destination")
		}

		json.NewEncoder(rw).Encode(mockResp)
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
