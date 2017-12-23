package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestNoti(t *testing.T) {
	if _, err := exec.LookPath("noti"); err != nil {
		t.Skip("noti binary missing from PATH")
	}

	t.Run("show version", func(t *testing.T) {
		data, err := exec.Command("noti", "--version").Output()
		if err != nil {
			t.Error(err)
		}
		out := string(data)

		if !strings.Contains(out, "noti version") {
			t.Error("Missing 'noti version'")
		}
		if !strings.Contains(out, "Latest:") {
			t.Error("Missing name of latest version")
		}
		if !strings.Contains(out, "Download: https://github.com/variadico/noti/releases") {
			t.Error("Missing latest download link")
		}
	})

	t.Run("dry run", func(t *testing.T) {
		cmd := exec.Command("noti", "--verbose", "-b=0")
		cmd.Env = []string{}

		data, err := cmd.Output()
		if err != nil {
			t.Error(err)
		}
		out := string(data)

		if !strings.Contains(out, "0 notifications queued") {
			t.Error("Unexpected queued notifications")
			t.Error(out)
		}
	})
}
