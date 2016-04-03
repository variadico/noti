// +build linux freebsd

package speech

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/variadico/noti"
)

// Notify speaks a notification using eSpeak.
func Notify(n noti.Params) error {
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

	voice := n.Config.Get(voiceEnv)
	if voice == "" {
		voice = "english-us"
	}
	text := fmt.Sprintf("%s %s", n.Title, n.Message)

	cmd := exec.Command("espeak", "-v", voice, "--", text)
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("Speech: %s", err)
	}

	return nil
}
