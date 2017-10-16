package bearychat

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	n := Notification{
		Text:            "mesg",
		IncomingHookURI: "",
		Client:          &http.Client{Timeout: 3 * time.Second},
	}
	var mockResp apiResponse
	var hitServer bool

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		hitServer = true

		if r.Method != "POST" {
			t.Error("HTTP method should be POST.")
		}
		defer r.Body.Close()

		var payload Notification
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("decoding incoming payload: %s", err)
		}
		if payload.Text != n.Text {
			t.Error("missing payload text")
		}

		json.NewEncoder(rw).Encode(mockResp)
	}))
	defer ts.Close()

	n.IncomingHookURI = ts.URL
	mockResp.Code = 0 // successful
	if err := n.Send(); err != nil {
		t.Error(err)
	}

	if !hitServer {
		t.Error("didn't reach server")
	}

	mockResp.Code = 1 // failure
	if err := n.Send(); err == nil {
		t.Error("unexpected success")
	}
}
