package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"
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

	t.Run("pwatch", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		sleep := exec.CommandContext(ctx, "sleep", "2")
		if err := sleep.Start(); err != nil {
			t.Error(err)
		}
		go sleep.Wait()

		var sleepPID string
		for i := 0; i < 3; i++ {
			if sleep.Process == nil {
				time.Sleep(10 * time.Millisecond)
				continue
			}
			sleepPID = fmt.Sprint(sleep.Process.Pid)
		}
		if sleepPID == "" {
			t.Fatal("couldn't get sleep pid")
		}

		cmd := exec.CommandContext(ctx, "noti", "--verbose", "-b=0", "--pwatch", sleepPID)
		cmd.Env = []string{}

		if out, err := cmd.Output(); err != nil {
			t.Errorf("noti: %s", err)
			t.Error(string(out))
		}
	})
}
