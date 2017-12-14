package hipchat

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
		Message: "mesg",
		Client:  &http.Client{Timeout: 3 * time.Second},
		Token:   "foo",
	}
	var mockResp apiResponse
	var hitServer bool

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		hitServer = true

		if r.Method != "POST" {
			t.Error("HTTP method should be POST")
		}

		if r.Header.Get("Authorization") == "" {
			t.Error("missing access token")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("content type should be application/json")
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}

		if string(b) == "" {
			t.Error("missing payload")
		}

		if mockResp.Error.Message != "" {
			json.NewEncoder(rw).Encode(mockResp)
		}
		// no response means success
	}))
	defer ts.Close()

	// HipChat API URL needs to be calculated based on the destination variable.
	API = ts.URL + "/%s"
	mockResp.Error.Message = "" // success
	if err := n.Send(); err != nil {
		t.Error(err)
	}

	if !hitServer {
		t.Error("Didn't reach server.")
	}

	mockResp.Error.Message = "error" // failure
	if err := n.Send(); err == nil {
		t.Error("unexpected success")
	}
}
