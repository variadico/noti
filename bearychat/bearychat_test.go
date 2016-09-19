package bearychat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/variadico/noti"
)

func TestEnvConfig(t *testing.T) {
	cases := []struct {
		env     noti.MockEnv
		config  configuration
		wantErr bool
		err     error
	}{
		{
			env:     noti.MockEnv{},
			config:  configuration{},
			wantErr: true,
			err:     errNoIncoming,
		},
		{
			env:     noti.MockEnv{incomingEnv: "foo"},
			config:  configuration{"foo"},
			wantErr: false,
			err:     nil,
		},
	}

	for i, c := range cases {
		config, err := envConfig(c.env)
		gotErr := (err != nil)

		if gotErr && !c.wantErr {
			t.Error(i, "unexpected error")
			t.Error(err)
		} else if gotErr && c.wantErr {
			if err != c.err {
				t.Error(i, "unexpected error")
				t.Error(err)
			}
		}

		if config != c.config {
			t.Error(i, "unexpected configuration")
			t.Errorf(" got: %v", config)
			t.Errorf("want: %v", c.config)
		}
	}
}

func TestNotify(t *testing.T) {
	n := noti.Params{
		Title:   "title",
		Message: "mesg",
		API:     "",
		Config:  nil,
	}
	var mockResp incomingResponse
	var hitServer bool

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		hitServer = true

		if r.Method != "POST" {
			t.Error("HTTP method should be POST.")
		}
		defer r.Body.Close()

		var payload incomingPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("decoding incoming payload: %s", err)
		}
		if payload.Text != fmt.Sprintf("**%s**\n%s", n.Title, n.Message) {
			t.Error("missing payload text")
		}

		json.NewEncoder(rw).Encode(mockResp)
	}))
	defer ts.Close()

	n.Config = noti.MockEnv{incomingEnv: ts.URL}
	mockResp.Code = 0 // successful
	if err := Notify(n); err != nil {
		t.Error(err)
	}

	if !hitServer {
		t.Error("didn't reach server")
	}

	mockResp.Code = 1 // failure
	if err := Notify(n); err == nil {
		t.Error("unexpected success")
	}
}
