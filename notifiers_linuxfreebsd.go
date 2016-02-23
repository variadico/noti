// +build linux freebsd

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

const (
	specificPart = `
    NOTI_VOICE
        Name of voice used for speech notifications. See "espeak --voices" for
        available voices.`
)

func init() {
	flag.Usage = func() {
		fmt.Printf(manual, specificPart)
	}
}

// bannerNotify triggers a Notify notification.
func bannerNotify() error {
	_, err := exec.LookPath("notify-send")
	if err != nil {
		return errors.New("Install 'notify-send' and try again")
	}

	cmd := exec.Command("notify-send", *title, *message)
	if err = cmd.Run(); err != nil {
		return err
	}

	return nil
}

// speechNotify triggers an eSpeak notification.
func speechNotify() error {
	_, err := exec.LookPath("espeak")
	if err != nil {
		buf := new(bytes.Buffer)
		buf.WriteString("Install 'espeak' and try again\n")

		var errStr string
		if runtime.GOOS == "freebsd" {
			buf.WriteString("On FreeBSD this might be: 'sudo pkg install --yes espeak'")
		} else if runtime.GOOS == "linux" {
			buf.WriteString("On Linux this might be: 'sudo apt-get install --yes espeak'")
		}

		return err
	}

	voice := os.Getenv(voiceEnv)
	if voice == "" {
		voice = "english-us"
	}

	*message = fmt.Sprintf("%s %s", *title, *message)

	cmd := exec.Command("espeak", "-v", voice, "--", *message)
	if err = cmd.Run(); err != nil {
		return err
	}

	return nil
}
