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
	script := buildAppleScript(n)
	cmd := exec.Command("osascript", "-e", script)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// escapeAppleScript escapes a string for use in an AppleScript double-quoted string.
func escapeAppleScript(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	return s
}

func buildAppleScript(n *Notification) string {
	var parts []string
	parts = append(parts, fmt.Sprintf(`display notification "%s"`, escapeAppleScript(n.InformativeText)))
	if n.Title != "" {
		parts = append(parts, fmt.Sprintf(`with title "%s"`, escapeAppleScript(n.Title)))
	}
	if n.Subtitle != "" {
		parts = append(parts, fmt.Sprintf(`subtitle "%s"`, escapeAppleScript(n.Subtitle)))
	}
	if n.SoundName != "" {
		parts = append(parts, fmt.Sprintf(`sound name "%s"`, escapeAppleScript(n.SoundName)))
	}
	return strings.Join(parts, " ")
}
