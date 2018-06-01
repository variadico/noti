package mattermost

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
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
			w.Write([]byte(`<HTML><HEAD><meta http-equiv="content-type" content="text/html;charset=utf-8"></HEAD><BODY></BODY></HTML>`))
		case "error":
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"id":"Unable to parse incoming data","message":"Unable to parse incoming data","detailed_error":"","request_id":"OoFaz2ra0Bahsiechu8c","status_code":400}`))
		default:

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
	if err := n.Send(); err.Error() != `response: {Unable to parse incoming data Unable to parse incoming data  %!s(int=400) OoFaz2ra0Bahsiechu8c}` {
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

// testMessage validate the payload the request method
// and if the neccessary field text is included
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
