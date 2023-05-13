package pushbullet

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	n := Notification{
		Title:       "title",
		Body:        "mesg",
		AccessToken: "token",
		Client:      &http.Client{Timeout: 3 * time.Second},
	}
	var mockResp apiResponse
	var hitServer bool

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		hitServer = true

		if r.Method != "POST" {
			t.Error("HTTP method should be POST")
		}

		if r.Header.Get("Access-Token") == "" {
			t.Error("missing access token")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("content type should be application/json")
		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}

		if string(b) == "" {
			t.Error("missing payload")
		}

		json.NewEncoder(rw).Encode(mockResp)
	}))
	defer ts.Close()

	API = ts.URL
	mockResp.ErrorCode = "" // success
	if err := n.Send(); err != nil {
		t.Error(err)
	}

	if !hitServer {
		t.Error("didn't reach server")
	}

	mockResp.ErrorCode = "error" // failure
	if err := n.Send(); err == nil {
		t.Error("unexpected success")
	}
}
