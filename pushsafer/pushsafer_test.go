package pushsafer

import (
	"encoding/json"
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
			err:     errNoKey,
		},
		{
			env:     noti.MockEnv{keyEnv: "fu"},
			config:  configuration{privateKey: "fu"},
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
		API:     API,
		Config:  noti.MockEnv{keyEnv: "fu"},
	}
	var mockResp apiResponse
	var hitServer bool

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		hitServer = true

		if r.Method != "POST" {
			t.Error("HTTP method should be POST")
		}

		if r.FormValue("k") == "" {
			t.Error("missing private or alias key")
		}

		json.NewEncoder(rw).Encode(mockResp)
	}))
	defer ts.Close()

	n.API = ts.URL
	mockResp.Status = 1 // success
	if err := Notify(n); err != nil {
		t.Error(err)
	}

	if !hitServer {
		t.Error("didn't reach server")
	}

	mockResp.Status = 0 // failure
	mockResp.Errors = []string{"error fu"}
	if err := Notify(n); err == nil {
		t.Error("unexpected success")
	}

	mockResp.Status = 1 // failure
	mockResp.Info = "no active devices to send to"
	if err := Notify(n); err == nil {
		t.Error("unexpected success")
	}
}
