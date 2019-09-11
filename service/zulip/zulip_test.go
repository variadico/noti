package zulip

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	n := Notification{
		BotAPIKey:       "key",
		BotEmailAddress: "from",
		Content:         "content",
		Type:            "stream",
		To:              "to",
		Endpoint:        "https://test/v1",
		Client:          &http.Client{Timeout: 3 * time.Second},
	}

	var mockResponse apiResponse
	var hitServer bool

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		hitServer = true

		if r.Method != "POST" {
			t.Error("HTTP method should be POST")
		}

		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Error("content type should be application/x-www-form-urlencoded")
		}

		user, pass, ok := r.BasicAuth()
		if user != "from" {
			t.Error("missing auth username")
		}
		if pass != "key" {
			t.Error("missing auth password")
		}
		if ok != true {
			t.Error("missing auth")
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}

		if string(b) == "" {
			t.Error("missing payload")
		}

		if err := json.NewEncoder(rw).Encode(mockResponse); err != nil {
			t.Error(err)
		}
	}))
	defer ts.Close()

	n.Endpoint = ts.URL
	mockResponse.Result = "success"
	if err := n.Send(); err != nil {
		t.Error(err)
	}

	if !hitServer {
		t.Error("didn't reach server")
	}

	n.Type = "private"
	if err := n.Send(); err != nil {
		t.Error(err)
	}

	// failure
	mockResponse.Result = "error"
	if err := n.Send(); err == nil {
		t.Error("unexpected success")
	}
	n.Type = "dummy"
	if err := n.Send(); err == nil {
		t.Error("unexpected success")
	}
}
