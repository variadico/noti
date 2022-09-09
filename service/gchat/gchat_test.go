package gchat

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	n := Notification{
		Title:    "title",
		Message:  "mesg",
		Client:   &http.Client{Timeout: 3 * time.Second},
		Template: "something something",
	}

	var hitServer bool
	var statusCode int

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		hitServer = true

		if r.Method != "POST" {
			t.Error("HTTP method should be POST")
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("content type should be application/json")
		}

		msg := &gchatMessage{}
		err := json.NewDecoder(r.Body).Decode(&msg)

		if err != nil {
			t.Error(err)
		}
		if statusCode == 0 {
			rw.WriteHeader(http.StatusOK)
		} else {
			rw.WriteHeader(statusCode)
			_, err := rw.Write([]byte("Error"))
			if err != nil {
				t.Error(err)
			}
		}
	}))
	defer ts.Close()

	n.AppURL = ts.URL
	if err := n.Send(); err != nil {
		t.Error(err)
	}

	if !hitServer {
		t.Error("didn't reach server")
	}

	statusCode = http.StatusNotImplemented
	if err := n.Send(); err == nil {
		t.Error("unexpected success")
	}
}
