package main

import (
	"flag"
	"os"
	"os/exec"
	"testing"
)

func TestNewNotification(t *testing.T) {
	orig := struct {
		runFn   func(*exec.Cmd) error
		flagSet *flag.FlagSet
	}{
		runFn:   run,
		flagSet: flag.CommandLine,
	}
	defer func() {
		run = orig.runFn
		flag.CommandLine = orig.flagSet
	}()
	run = func(*exec.Cmd) error {
		return nil
	}

	type tc []struct {
		input []string
		want  notification
	}
	var tests tc

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flag.CommandLine.Parse([]string{})
	tests = tc{
		{[]string{}, notification{"noti", "Done!", false}},
		{[]string{"ls"}, notification{"ls", "Done!", false}},
		{[]string{"ls", "-al"}, notification{"ls", "Done!", false}},
	}
	for i, test := range tests {
		n := newNotification(test.input)
		if n.title != test.want.title {
			t.Errorf("%d: got: %q; want: %q\n", i, n.title, test.want.title)
		}
		if n.message != test.want.message {
			t.Errorf("%d: got: %q; want: %q\n", i, n.message, test.want.message)
		}
	}

	tests = tc{
		{[]string{}, notification{"fu", "Done!", false}},
		{[]string{"ls"}, notification{"fu", "Done!", false}},
		{[]string{"ls", "-al"}, notification{"fu", "Done!", false}},
	}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	title = flag.String("t", "", "")
	flag.CommandLine.Parse([]string{"-t", "fu"})
	for i, test := range tests {
		n := newNotification(test.input)
		if n.title != test.want.title {
			t.Errorf("%d: got: %q; want: %q\n", i, n.title, test.want.title)
		}
		if n.message != test.want.message {
			t.Errorf("%d: got: %q; want: %q\n", i, n.message, test.want.message)
		}
	}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	title = flag.String("title", "", "")
	flag.CommandLine.Parse([]string{"-title", "fu"})
	for i, test := range tests {
		n := newNotification(test.input)
		if n.title != test.want.title {
			t.Errorf("%d: got: %q; want: %q\n", i, n.title, test.want.title)
		}
		if n.message != test.want.message {
			t.Errorf("%d: got: %q; want: %q\n", i, n.message, test.want.message)
		}
	}

	tests = tc{
		{[]string{}, notification{"noti", "fa", false}},
		{[]string{"ls"}, notification{"ls", "fa", false}},
		{[]string{"ls", "-al"}, notification{"ls", "fa", false}},
	}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	message = flag.String("m", "", "")
	flag.CommandLine.Parse([]string{"-m", "fa"})
	for i, test := range tests {
		n := newNotification(test.input)
		if n.title != test.want.title {
			t.Errorf("%d: got: %q; want: %q\n", i, n.title, test.want.title)
		}
		if n.message != test.want.message {
			t.Errorf("%d: got: %q; want: %q\n", i, n.message, test.want.message)
		}
	}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	message = flag.String("message", "", "")
	flag.CommandLine.Parse([]string{"-message", "fa"})
	for i, test := range tests {
		n := newNotification(test.input)
		if n.title != test.want.title {
			t.Errorf("%d: got: %q; want: %q\n", i, n.title, test.want.title)
		}
		if n.message != test.want.message {
			t.Errorf("%d: got: %q; want: %q\n", i, n.message, test.want.message)
		}
	}
}

func TestRunUtility(t *testing.T) {
	orig := run
	defer func() {
		run = orig
	}()

	var hitRun bool
	run = func(*exec.Cmd) error {
		hitRun = true
		return nil
	}

	util, err := runUtility([]string{})
	if err != nil {
		t.Error("Empty slice shouldn't cause error.")
	}
	if util != "noti" {
		t.Error("Empty slice should default util to 'noti'.")
	}
	if hitRun {
		t.Error("Empty slice shouldn't try to execute anything.")
	}

	util, err = runUtility([]string{"ls"})
	if err != nil {
		t.Error(err)
	}
	if util[:2] != "ls" {
		t.Error("Utility name should be returned.")
	}
}
