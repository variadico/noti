// +build darwin

package speech

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/variadico/noti"
)

// speechNotify triggers an NSSpeechSynthesizer notification.
func Notify(n noti.Notification) error {
	voice := n.Config.Get(voiceEnv)
	if voice == "" {
		voice = "Alex"
	}
	text := fmt.Sprintf("%s %s", n.Title, n.Message)

	cmd := exec.Command("say", "-v", voice, text)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Speech: %s", err)
	}

	return nil
}
