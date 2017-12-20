package command

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestRunCommand(t *testing.T) {
	lsAbs, err := exec.LookPath("ls")
	if err != nil {
		if runtime.GOOS == "windows" {
			lsAbs, err = exec.LookPath("dir")
		}
		t.Skipf("Missing command for test, skipping: %s", err)
	}

	if err := runCommand([]string{lsAbs}, nil, nil, nil); err != nil {
		t.Error("Unexpected error", err)
	}

	if err := runCommand([]string{}, nil, nil, nil); err != nil {
		t.Error("Unexpected error", err)
	}

	shell := filepath.Base(os.Getenv("SHELL"))
	if shell != "bash" && shell != "zsh" {
		t.Logf("%s subshell not supported", shell)
		t.Log("Skipping rest of test")
		return
	}

	if err := runCommand([]string{":"}, nil, nil, nil); err != nil {
		t.Error("Unexpected error", err)
	}
}

func TestCommandName(t *testing.T) {
	cases := []struct {
		args []string
		want string
	}{
		{args: []string{},
			want: "noti"},
		{args: []string{"git"},
			want: "git"},
		{args: []string{"git", "commit"},
			want: "git commit"},
		{args: []string{"ls", "-l"},
			want: "ls"},
		{args: []string{"foo", "bar", "fizz", "buzz"},
			want: "foo bar"},
	}

	for i, c := range cases {
		have := commandName(c.args)
		if have != c.want {
			t.Error("Unexpected command name")
			t.Errorf("%d - have=%q; want=%q", i, have, c.want)
		}
	}
}

func TestLatestRelease(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Error("Unexpected request method")
			t.Errorf("have=%s; want=%s", r.Method, "GET")
		}

		data, err := ioutil.ReadFile("testdata/github_latest_release.json")
		if err != nil {
			t.Error(err)
		}
		fmt.Fprintln(rw, string(data))
	}))
	defer ts.Close()

	haveLatest, haveDownload, err := latestRelease(ts.URL)
	if err != nil {
		t.Error(err)
	}

	const wantLatest = "v1.0.0"
	if haveLatest != wantLatest {
		t.Error("Unexpected latest tag")
		t.Errorf("have=%s; want=%s", haveLatest, wantLatest)
	}

	const wantDownload = "https://github.com/octocat/Hello-World/releases/v1.0.0"
	if haveDownload != wantDownload {
		t.Error("Unexpected download URL")
		t.Errorf("have=%s; want=%s", haveDownload, wantDownload)
	}
}
