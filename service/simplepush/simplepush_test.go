package simplepush

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	n := Notification{
		Title:   "title",
		Message: "mesg",
		Key:     "key",
		Client:  &http.Client{Timeout: 3 * time.Second},
	}
	var mockResp apiResponse
	var hitServer bool

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		hitServer = true

		if r.Method != "POST" {
			t.Error("HTTP method should be POST")
		}

		if r.FormValue("key") == "" {
			t.Error("missing key")
		}

		json.NewEncoder(rw).Encode(mockResp)
	}))
	defer ts.Close()

	API = ts.URL
	mockResp.Status = "OK" // success
	if err := n.Send(); err != nil {
		t.Error(err)
	}

	if !hitServer {
		t.Error("didn't reach server")
	}

	mockResp.Status = "BadRequest" // failure
	if err := n.Send(); err == nil {
		t.Error("unexpected success")
	}
}
