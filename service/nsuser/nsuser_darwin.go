package nsuser

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Notification is a macOS banner notification.
type Notification struct {
	Title    string
	Subtitle string
	// InformativeText is the notification message.
	InformativeText string
	// ContentImage is unused but kept for API compatibility.
	ContentImage string
	// SoundName is the name of the sound that fires with a notification.
	SoundName string
}

// Send displays a macOS banner notification via osascript.
func (n *Notification) Send() error {
	if n.ContentImage != "" {
		log.Println("nsuser: ContentImage is not supported with osascript notifications; ignoring")
	}
	script := buildScript(n)
	cmd := exec.Command("osascript", "-e", script)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// escapeAS escapes a string for use in an AppleScript double-quoted string.
func escapeAS(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	return s
}

func buildScript(n *Notification) string {
	var parts []string
	parts = append(parts, fmt.Sprintf(`display notification "%s"`, escapeAS(n.InformativeText)))
	if n.Title != "" {
		parts = append(parts, fmt.Sprintf(`with title "%s"`, escapeAS(n.Title)))
	}
	if n.Subtitle != "" {
		parts = append(parts, fmt.Sprintf(`subtitle "%s"`, escapeAS(n.Subtitle)))
	}
	if n.SoundName != "" {
		parts = append(parts, fmt.Sprintf(`sound name "%s"`, escapeAS(n.SoundName)))
	}
	return strings.Join(parts, " ")
}
