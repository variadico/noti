package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/variadico/noti"
	"github.com/variadico/noti/banner"
	"github.com/variadico/noti/bearychat"
	"github.com/variadico/noti/hipchat"
	"github.com/variadico/noti/pushbullet"
	"github.com/variadico/noti/pushover"
	"github.com/variadico/noti/slack"
	"github.com/variadico/noti/speech"
)

const (
	defaultEnv = "NOTI_DEFAULT"
	version    = "v2.3.0"
)

func main() {
	log.SetFlags(0)

	title := flag.String("t", "", "")
	message := flag.String("m", "", "")
	showVersion := flag.Bool("v", false, "")
	showHelp := flag.Bool("h", false, "")
	pid := flag.Int("w", -1, "")

	// Notifications
	bannerNoti := flag.Bool("b", true, "")
	hipChatNoti := flag.Bool("i", false, "")
	pushbulletNoti := flag.Bool("p", false, "")
	pushoverNoti := flag.Bool("o", false, "")
	slackNoti := flag.Bool("k", false, "")
	speechNoti := flag.Bool("s", false, "")
	bearychatNoti := flag.Bool("c", false, "")

	flag.StringVar(title, "title", "", "")
	flag.StringVar(message, "message", "", "")
	flag.IntVar(pid, "pwatch", -1, "")
	flag.BoolVar(showVersion, "version", false, "")
	flag.BoolVar(showHelp, "help", false, "")

	// Notifications
	flag.BoolVar(bannerNoti, "banner", true, "")
	flag.BoolVar(hipChatNoti, "hipchat", false, "")
	flag.BoolVar(speechNoti, "speech", false, "")
	flag.BoolVar(pushbulletNoti, "pushbullet", false, "")
	flag.BoolVar(pushoverNoti, "pushover", false, "")
	flag.BoolVar(slackNoti, "slack", false, "")
	flag.BoolVar(bearychatNoti, "bearychat", false, "")

	flag.Parse()

	if *showVersion {
		fmt.Println("noti", version)
		checkForUpdates()
		return
	}
	if *showHelp {
		flag.Usage()
		return
	}

	env := noti.OSEnv{}
	setDefaultNotifications(flag.CommandLine, env)

	var n noti.Params
	if *pid == -1 {
		n = utilityNotification(flag.CommandLine)
	} else {
		n = pwatchNotification(flag.CommandLine, *pid)
	}

	n.Config = env

	notis := []struct {
		trigger bool
		api     string
		notify  func(noti.Params) error
	}{
		{*bannerNoti, "", banner.Notify},
		{*hipChatNoti, hipchat.API, hipchat.Notify},
		{*pushbulletNoti, pushbullet.API, pushbullet.Notify},
		{*pushoverNoti, pushover.API, pushover.Notify},
		{*speechNoti, "", speech.Notify},
		{*slackNoti, slack.API, slack.Notify},
		{*bearychatNoti, "", bearychat.Notify},
	}

	for _, nt := range notis {
		if !nt.trigger {
			continue
		}

		n.API = nt.api

		if err := nt.notify(n); err != nil {
			log.Println(err)
		}
	}

	if n.Failure {
		os.Exit(1)
	}
}

func checkForUpdates() error {
	// Draft releases and prereleases are not returned by this endpoint.
	const releaseAPI = "https://api.github.com/repos/variadico/noti/releases/latest"
	webClient := &http.Client{Timeout: 30 * time.Second}

	resp, err := webClient.Get(releaseAPI)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r struct {
		HTMLURL string `json:"html_url"`
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	fmt.Println("latest:", r.TagName)
	if r.TagName != version {
		fmt.Println("download:", r.HTMLURL)
	}

	return nil
}

// setDefaultNotifications read the user's config and set their defaults on a
// FlagSet.
func setDefaultNotifications(fl *flag.FlagSet, env noti.EnvGetter) {
	defs := strings.TrimSpace(env.Get(defaultEnv))
	if defs == "" {
		return
	}

	has := strings.Contains
	fl.Set("banner", fmt.Sprintf("%t", has(defs, "banner")))
	fl.Set("hipchat", fmt.Sprintf("%t", has(defs, "hipchat")))
	fl.Set("pushbullet", fmt.Sprintf("%t", has(defs, "pushbullet")))
	fl.Set("pushover", fmt.Sprintf("%t", has(defs, "pushover")))
	fl.Set("speech", fmt.Sprintf("%t", has(defs, "speech")))
	fl.Set("slack", fmt.Sprintf("%t", has(defs, "slack")))
}

func pwatchNotification(fl *flag.FlagSet, pid int) noti.Params {
	err := pollPID(pid, 1*time.Second)

	return noti.Params{
		Title:   notiTitle(fl, "noti", err),
		Message: notiMessage(fl, err),
		Failure: (err != nil),
	}
}

func utilityNotification(fl *flag.FlagSet) noti.Params {
	util, err := runUtility(fl.Args())

	return noti.Params{
		Title:   notiTitle(fl, util, err),
		Message: notiMessage(fl, err),
		Failure: (err != nil),
	}
}

func notiTitle(fl *flag.FlagSet, util string, err error) string {
	t := flagValue(fl, "t", "title")
	if t == "" {
		t = util
	}

	if err != nil {
		t = fmt.Sprintf("%s failed", t)
	}

	return t
}

func notiMessage(fl *flag.FlagSet, err error) string {
	m := flagValue(fl, "m", "message")
	if m == "" {
		m = "Done!"
	}

	if err != nil {
		m = err.Error()
	}

	return m
}

func flagValue(fl *flag.FlagSet, short, long string) string {
	var v string

	if userSet(fl, short) {
		v = fl.Lookup(short).Value.String()
	} else if userSet(fl, long) {
		v = fl.Lookup(long).Value.String()
	}

	return v
}

func runUtility(args []string) (string, error) {
	name := utilityName(args)
	if len(args) < 1 {
		return name, nil
	}

	if _, err := exec.LookPath(args[0]); err != nil {
		exp, expErr := expandAlias(args[0])
		if expErr != nil {
			return "", err
		}

		args = append(exp, args[1:]...)
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return name, cmd.Run()
}

func utilityName(args []string) string {
	switch len(args) {
	case 0:
		return "noti"
	case 1:
		return args[0]
	}

	// If the next arg isn't a flag.
	if args[1][0] != '-' {
		// Append a subcommand to the utility name.
		return fmt.Sprintf("%s %s", args[0], args[1])
	}

	return args[0]
}

// userSet returns true if a user passed a target flag. Otherwise, it returns
// false. With zero values, we can't easily tell if a user omited a flag or if
// she actually did passed it with a zero value.
func userSet(fl *flag.FlagSet, target string) bool {
	var explicitlySet bool

	fl.Visit(func(f *flag.Flag) {
		if f.Name == target {
			explicitlySet = true
		}
	})

	return explicitlySet
}

func expandAlias(a string) ([]string, error) {
	shell := os.Getenv("SHELL")

	cmd := exec.Command(shell, "-l", "-i", "-c", "which "+a)
	e, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	exp := strings.TrimSpace(string(e))
	trimLen := fmt.Sprintf("%s: aliased to ", a)
	exp = exp[len(trimLen):]

	return strings.Split(exp, " "), nil
}
