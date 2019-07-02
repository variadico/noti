package telegram

import "net/http"

type apiResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		MessageID int `json:"message_id"`
		Chat      struct {
			ID       uint8  `json:"id"`
			Title    string `json:"title"`
			Username string `json:"username"`
			Type     string `json:"type"`
		} `json:"chat"`
		Date int16  `json:"date"`
		Text string `json:"text"`
	} `json:"result"`
}

type Notification struct {
	ChatID  string
	Message string
	Client  *http.Client
}

func (n *Notification) Send() error {
	return nil
}
