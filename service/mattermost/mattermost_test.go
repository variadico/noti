package mattermost

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	n := Notification{
		Text:            "msg",
		IncomingHookURI: "",
		Client:          &http.Client{Timeout: 3 * time.Second},
	}
	//	var mockResp apiErrorResponse
	var hitServer bool

	// Setup testing Suite
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			t.Error("Could not Parse Query parameters")
		}
		hitServer = true
		checkMessage(r, n, t)

		action := r.Form.Get("action")
		switch action {
		case "good":
			w.Write([]byte("ok"))
		case "bad":
			w.Write([]byte(
				`<HTML><HEAD>
				<meta http-equiv="content-type" content="text/html;charset=utf-8">
				</HEAD> <BODY></BODY></HTML>`,
			))
		case "error":
			w.WriteHeader(http.StatusBadRequest)
			data := struct {
				ID            string `json:"id"`
				Message       string `json:"message"`
				DetailedError string `json:"detailed_error"`
				RequestID     string `json:"request_id"`
				StatusCode    int    `json:"status_code"`
			}{
				ID:         "Unable to parse incoming data",
				Message:    "Unable to parse incoming data",
				RequestID:  "OoFaz2ra0Bahsiechu8c",
				StatusCode: http.StatusBadRequest,
			}

			if err := json.NewEncoder(w).Encode(data); err != nil {
				t.Error(err)
			}
		}

	}))
	defer ts.Close()

	hook, err := url.Parse(ts.URL)
	if err != nil {
		t.Error("Could not Parse URL")
	}

	// check good
	hook.RawQuery = GetQueryValues("action=good", t).Encode()
	n.IncomingHookURI = hook.String()
	if err := n.Send(); err != nil {
		t.Error(err)
	}
	if !hitServer {
		t.Error("didn't reach server")
	}

	// Check error
	hook.RawQuery = GetQueryValues("action=error", t).Encode()
	n.IncomingHookURI = hook.String()
	if err := n.Send(); !strings.Contains(err.Error(), "Unable to parse incoming data") {
		t.Error(err)
	}

	if !hitServer {
		t.Error("didn't reach server")
	}

	// Check bad
	hook.RawQuery = GetQueryValues("action=bad", t).Encode()
	n.IncomingHookURI = hook.String()
	if err := n.Send(); err != ErrInvalidResponse {
		t.Error(err)
	}

	if !hitServer {
		t.Error("didn't reach server")
	}
}

func GetQueryValues(rawStr string, t *testing.T) url.Values {
	query, err := url.ParseQuery(rawStr)
	if err != nil {
		t.Error(err)
	}
	return query
}

func checkMessage(r *http.Request, n Notification, t *testing.T) {
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
}
