package notifyicon

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

// Balloon icons.
const (
	BalloonTipIconError   = "Error"
	BalloonTipIconWarning = "Warning"
	BalloonTipIconInfo    = "Info"
	BalloonTipIconNone    = "None"

	// DefaultIcon is the default icon.
	DefaultIcon = "[System.Drawing.Icon]::ExtractAssociatedIcon([System.Windows.Forms.Application]::ExecutablePath)"
)

const script = `
[void] [System.Reflection.Assembly]::LoadWithPartialName("System.Windows.Forms")

$n = New-Object System.Windows.Forms.NotifyIcon
$n.Icon = {{.Icon}}
$n.BalloonTipIcon = "{{.BalloonTipIcon}}"
$n.BalloonTipText = "{{.BalloonTipText}}"
$n.BalloonTipTitle = "{{.BalloonTipTitle}}"
$n.Text = "{{.Text}}"

$n.Visible = $True
$n.ShowBalloonTip({{.Duration}})
`

// Notification is a Windows notification.
type Notification struct {
	// BalloonTipIcon is the notification icon.
	BalloonTipIcon string
	// BalloonTipText is the notification message.
	BalloonTipText string
	// BalloonTipTitle is the notification title.
	BalloonTipTitle string
	// Icon is the path to an .ico file.
	// Icon sets the icon that will appear in the systray for this application.
	// Icon must be 16 pixels high by 16 pixels wide
	// This is required to show the notification.
	Icon string

	// Text is the text shown when you hover over the app icon.
	Text string

	Duration int
}

// Send sends a Windows notification.
func (n *Notification) Send() error {
	if n.Icon == "" {
		n.Icon = DefaultIcon
	} else {
		n.Icon = fmt.Sprintf("%q", n.Icon)
	}

	tmpl, err := template.New("").Parse(script)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, n); err != nil {
		return err
	}

	cmd := exec.Command("PowerShell", "-Command", buf.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// SetMessage sets this notification's message.
func (n *Notification) SetMessage(m string) {
	n.BalloonTipText = m
}

// GetMessage gets this notification's message.
func (n *Notification) GetMessage() string {
	return n.BalloonTipText
}
