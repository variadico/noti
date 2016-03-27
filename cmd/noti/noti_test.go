package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/variadico/noti"
)

func TestSetDefaultNotifications(t *testing.T) {
	fl := &flag.FlagSet{Usage: func() {}}
	fl.SetOutput(ioutil.Discard)
	fl.Bool("banner", false, "")
	fl.Bool("hipchat", false, "")
	fl.Bool("pushbullet", false, "")
	fl.Bool("pushover", false, "")
	fl.Bool("speech", false, "")
	fl.Bool("slack", false, "")
	if err := fl.Parse([]string{}); err != nil {
		t.Fatal(err)
	}

	notis := []string{
		"banner",
		"hipchat",
		"pushbullet",
		"pushover",
		"speech",
		"slack",
	}
	mockEnv := noti.MockEnv{}

	// Without config, don't touch flag set.
	setDefaultNotifications(fl, mockEnv)
	for i, nt := range notis {
		enabled := fl.Lookup(nt).Value.String()
		if enabled != "false" {
			t.Error(i, "unexpected result")
			t.Errorf(" got: %s", enabled)
			t.Errorf("want: %s", "false")
		}
	}

	// With config, flag set should match confing.
	mockEnv[defaultEnv] = strings.Join(notis, " ")
	setDefaultNotifications(fl, mockEnv)
	for i, nt := range notis {
		enabled := fl.Lookup(nt).Value.String()
		if enabled != "true" {
			t.Error(i, "unexpected result")
			t.Errorf(" got: %s", enabled)
			t.Errorf("want: %s", "true")
		}
	}
}

func TestNewNotification(t *testing.T) {
	fl := &flag.FlagSet{Usage: func() {}}
	fl.SetOutput(ioutil.Discard)
	if err := fl.Parse([]string{}); err != nil {
		t.Fatal(err)
	}

	n := newNotification(fl)
	want := noti.Notification{
		Title:   "noti",
		Message: "Done!",
		Failure: false,
	}

	if n != want {
		t.Error("unexpected result")
		t.Errorf(" got: %#v", n)
		t.Errorf("want: %#v", want)
	}
}

func TestNotificationMessage(t *testing.T) {
	cases := []struct {
		flags *flag.FlagSet
		args  []string
		err   error
		want  string
	}{
		{
			flags: &flag.FlagSet{Usage: func() {}},
			args:  []string{"-m", "fu"},
			want:  "fu",
		},
		{
			flags: &flag.FlagSet{Usage: func() {}},
			args:  []string{"-message", "fu"},
			want:  "fu",
		},
		{
			flags: &flag.FlagSet{Usage: func() {}},
			args:  []string{},
			want:  "Done!",
		},
		{
			flags: &flag.FlagSet{Usage: func() {}},
			args:  []string{"-m", "fu", "ls"},
			err:   errors.New("error fa"),
			want:  "error fa",
		},
	}

	for i, c := range cases {
		c.flags.SetOutput(ioutil.Discard)
		c.flags.String("m", "", "")
		c.flags.String("message", "", "")
		if err := c.flags.Parse(c.args); err != nil {
			t.Fatal(err)
		}

		title := notificationMessage(c.flags, c.err)
		if title != c.want {
			t.Error(i, "unexpected result")
			t.Errorf(" got: %s", title)
			t.Errorf("want: %s", c.want)
		}
	}
}

func TestNotificationTitle(t *testing.T) {
	cases := []struct {
		flags *flag.FlagSet
		args  []string
		util  string
		err   error
		want  string
	}{
		{
			flags: &flag.FlagSet{Usage: func() {}},
			args:  []string{"-t", "fu"},
			util:  "noti",
			want:  "fu",
		},
		{
			flags: &flag.FlagSet{Usage: func() {}},
			args:  []string{"-title", "fu"},
			util:  "noti",
			want:  "fu",
		},
		{
			flags: &flag.FlagSet{Usage: func() {}},
			args:  []string{},
			util:  "noti",
			want:  "noti",
		},
		{
			flags: &flag.FlagSet{Usage: func() {}},
			args:  []string{"-t", "fu", "ls"},
			util:  "ls",
			want:  "fu",
		},
		{
			flags: &flag.FlagSet{Usage: func() {}},
			args:  []string{"-t", "fu", "ls"},
			util:  "ls",
			err:   errors.New("fu"),
			want:  "fu failed",
		},
	}

	for i, c := range cases {
		c.flags.SetOutput(ioutil.Discard)
		c.flags.String("t", "", "")
		c.flags.String("title", "", "")
		if err := c.flags.Parse(c.args); err != nil {
			t.Fatal(err)
		}

		title := notificationTitle(c.flags, c.util, c.err)
		if title != c.want {
			t.Error(i, "unexpected result")
			t.Errorf(" got: %s", title)
			t.Errorf("want: %s", c.want)
		}
	}
}

func TestUtilityName(t *testing.T) {
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
		{args: []string{"fu", "fa", "ton", "son"},
			want: "fu fa"},
	}

	for i, c := range cases {
		name := utilityName(c.args)
		if name != c.want {
			t.Error(i, "unexpected result")
			t.Errorf(" got: %s", name)
			t.Errorf("want: %s", c.want)
		}
	}
}

func TestUserSet(t *testing.T) {
	cases := []struct {
		flags   *flag.FlagSet
		args    []string
		target  string
		wantSet bool
	}{
		{
			flags:   &flag.FlagSet{Usage: func() {}},
			args:    []string{"-fu", "-fa"},
			target:  "fu",
			wantSet: true,
		},
		{
			flags:   &flag.FlagSet{Usage: func() {}},
			args:    []string{"-fu"},
			target:  "fa",
			wantSet: false,
		},
	}

	for i, c := range cases {
		c.flags.SetOutput(ioutil.Discard)
		c.flags.Bool("fu", false, "")
		c.flags.Bool("fa", false, "")
		if err := c.flags.Parse(c.args); err != nil {
			t.Fatal(err)
		}

		set := userSet(c.flags, c.target)
		if set != c.wantSet {
			t.Error(i, "unexpected result")
			t.Errorf(" got: %t", set)
			t.Errorf("want: %t", c.wantSet)
		}
	}
}
