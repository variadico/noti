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
func bannerNotify(n notification) error {
	_, err := exec.LookPath("notify-send")
	if err != nil {
		return errors.New("Install 'notify-send' and try again")
	}

	cmd := exec.Command("notify-send", n.title, n.message)
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("Banner: %s", err)
	}

	return nil
}

// speechNotify triggers an eSpeak notification.
func speechNotify(n notification) error {
	_, err := exec.LookPath("espeak")
	if err != nil {
		buf := new(bytes.Buffer)
		buf.WriteString("Install 'espeak' and try again\n")

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
	text := fmt.Sprintf("%s %s", n.title, n.message)

	cmd := exec.Command("espeak", "-v", voice, "--", text)
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("Speech: %s", err)
	}

	return nil
}
