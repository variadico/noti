package chanify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	n := Notification{
		Title:             "title",
		Text:              "mesg",
		Client:            &http.Client{Timeout: 3 * time.Second},
		Sound:             true,
		Priority:          10,
		InterruptionLevel: "active",
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

		msg := &chanifyMessage{}
		err := json.NewDecoder(r.Body).Decode(&msg)

		if err != nil {
			t.Error(err)
		}

		if msg.Title != n.Title {
			t.Errorf("Expected Title to be %s but got %s", n.Title, msg.Title)
		}

		if msg.Text != n.Text {
			t.Errorf("Expected Text to be %s but got %s", n.Text, msg.Text)
		}

		if msg.Sound != n.Sound {
			t.Errorf("Expected Sound to be %v but got %v", n.Sound, msg.Sound)
		}

		if msg.Priority != n.Priority {
			t.Errorf("Expected Priority to be %d but got %d", n.Priority, msg.Priority)
		}

		if msg.InterruptionLevel != n.InterruptionLevel {
			t.Errorf("Expected InterruptionLevel to be %s but got %s", n.InterruptionLevel, msg.InterruptionLevel)
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

	n.ChannelURL = ts.URL
	if err := n.Send(); err != nil {
		t.Error(err)
	}

	if !hitServer {
		t.Error("didn't reach server")
	}

	statusCode = http.StatusInternalServerError
	if err := n.Send(); err == nil {
		t.Error("unexpected success")
	}
}
