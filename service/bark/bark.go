package bark

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Notification is a Bark notification.
type Notification struct {
	URL       string `json:"api_url"`
	DeviceKey string `json:"device_key"`
	Body      string `json:"body"`
	Title     string `json:"title"`

	Client *http.Client `json:"-"`
}

// Send sends a Bark notification.
func (n *Notification) Send() error {
	if n.DeviceKey == "" {
		return errors.New("missing DeviceKey")
	}

	jsonData, err := json.Marshal(n)
	if err != nil {
		return err
	}

	resp, err := n.Client.Post(n.URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("bark: error reading response %w", err)
		}
		return fmt.Errorf("bark: %d\n %s", resp.StatusCode, string(b))
	}

	return nil

}
