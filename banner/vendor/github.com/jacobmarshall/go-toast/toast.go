package go_toast

import (
    "github.com/nu7hatch/gouuid"
    "text/template"
    "bytes"
    "io/ioutil"
    "os"
    "os/exec"
    "path/filepath"
)

var toastTemplate *template.Template

func init() {
    toastTemplate = template.New("toast")
    toastTemplate.Parse(`
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.UI.Notifications.ToastNotification, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType = WindowsRuntime] | Out-Null

$APP_ID = '{{if .AppID}}{{.AppID}}{{else}}io.github.jacobmarshall.go-toast{{end}}'

$template = @"
<toast>
    <visual>
        <binding template="ToastGeneric">
            {{if .Icon}}
            <image placement="appLogoOverride" src="{{.Icon}}" />
            {{end}}
            {{if .Title}}
            <text>{{.Title}}</text>
            {{end}}
            {{if .Message}}
            <text>{{.Message}}</text>
            {{end}}
        </binding>
    </visual>
    {{if .Actions}}
    <actions>
        {{range .Actions}}
        <action activationType="{{.Type}}" content="{{.Label}}" arguments="{{.Arguments}}" />
        {{end}}
    </actions>
    {{end}}
</toast>
"@

$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
$xml.LoadXml($template)
$toast = New-Object Windows.UI.Notifications.ToastNotification $xml
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier($APP_ID).Show($toast)
    `);
}

type Notification struct {
    AppID string
    Title string
    Message string
    Icon string
    Actions []Action
}

type Action struct {
    Type string
    Label string
    Arguments string
}

func (n *Notification) BuildXML() (string, error) {
    var out bytes.Buffer
    err := toastTemplate.Execute(&out, n)
    if err != nil {
        return "", err
    }
    return out.String(), nil
}

func (n *Notification) Push() error {
    xml, _ := n.BuildXML()
    return invokeTemporaryScript(xml)
}

func invokeTemporaryScript(content string) (error) {
    id, _ := uuid.NewV4()
    file := filepath.Join(os.TempDir(), id.String() + ".ps1")
    defer os.Remove(file)
    err := ioutil.WriteFile(file, []byte(content), 0600)
    if err != nil {
        return err
    }
    if err = exec.Command("PowerShell", "-File", file).Run(); err != nil {
        return err
    }
    return nil
}
