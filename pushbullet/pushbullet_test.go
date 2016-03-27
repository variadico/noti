package pushbullet

import (
	"encoding/json"
	"io/ioutil"
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
			err:     errNoTok,
		},
		{
			env:     noti.MockEnv{tokEnv: "fu"},
			config:  configuration{"fu"},
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
	n := noti.Notification{
		Title:   "title",
		Message: "mesg",
		API:     API,
		Config:  noti.MockEnv{tokEnv: "fu"},
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

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}

		if string(b) == "" {
			t.Error("missing payload")
		}

		json.NewEncoder(rw).Encode(mockResp)
	}))
	defer ts.Close()

	n.API = ts.URL
	mockResp.ErrorCode = "" // success
	if err := Notify(n); err != nil {
		t.Error(err)
	}

	if !hitServer {
		t.Error("didn't reach server")
	}

	mockResp.ErrorCode = "fu" // failure
	if err := Notify(n); err == nil {
		t.Error("unexpected success")
	}
}
