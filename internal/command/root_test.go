package command

import (
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
