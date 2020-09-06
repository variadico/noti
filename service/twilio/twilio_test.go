package twilio

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestTwilioSend(t *testing.T) {
	n := Notification{
		Content:    "Sleep:Testing notification",
		NumberTo:   "+112341234",
		NumberFrom: "+9728848438",
		AccountSid: "sid",
		AuthToken:  "token",
	}

	var mockResp apiResponse
	var hitServer bool

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		hitServer = true

		body, err1 := ioutil.ReadAll(r.Body)
		if err1 != nil {
			log.Printf("Error reading body: %v", err1)
			return
		}
		var v url.Values
		v, err := url.ParseQuery(string(body))
		if err != nil {
			t.Error("Invalid query")
		}

		if r.Method != "POST" {
			t.Error("HTTP method should be POST.")
		}

		if v["To"][0] == "" {
			t.Error("invalid receiver number")
		}

		if v["From"][0] == "" {
			t.Error("invalid sender")
		}

		if v["Body"][0] == "" {
			t.Error("invalid msg")
		}

		if err := json.NewEncoder(rw).Encode(mockResp); err != nil {
			t.Error(err)
		}
	}))
	defer ts.Close()

	API = ts.URL

	// successful
	mockResp.ErrorCode = ""
	if err := n.Send(); err != nil {
		t.Error(err)
	}

	if !hitServer {
		t.Error("didn't reach server")
	}

	// failure
	mockResp.ErrorCode = "failed"
	if err := n.Send(); err == nil {
		t.Error("unexpected success")
	}
}
