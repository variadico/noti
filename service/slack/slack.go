package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const (
	// ParseNone autoformats the text in messages less.
	ParseNone = "none"
	// ParseFull autoformats message text more, like creating hyperlinks
	// automatically.
	ParseFull = "full"

	// LinkNamesOn enables making usernames hyperlinks.
	LinkNamesOn = 1
	// LinkNamesOff disables making usernames hyperlinks.
	LinkNamesOff = 0
)

// API is the endpoint that handles posting a message.
var API = "https://slack.com/api/chat.postMessage"

type apiResponse struct {
	OK        bool   `json:"ok"`
	Channel   string `json:"channel"`
	Timestamp string `json:"ts"`
	Message   struct {
		Text     string `json:"text"`
		Username string `json:"username"`
		Icons    struct {
			Emoji   string `json:"emoji"`
			Image64 string `json:"image_64"`
		} `json:"icons"`
		Type      string `json:"type"`
		Subtype   string `json:"subtype"`
		Timestamp string `json:"ts"`
	} `json:"message"`
	Error string `json:"error"`
}

// Notification is a Slack notification.
type Notification struct {
	// AppURL is your Slack App's webhook URL
	AppURL string
	// Token is a user's authentication token.
	Token string
	// Channel is a notification's destination. It can be a channel, private
	// group, or username.
	Channel string
	// Text is the notification's message.
	Text string
	// Parse is the mode used to parse text.
	Parse string
	// LinkNames converts usernames into links.
	LinkNames int
	// Attachments are rich text snippets.
	Attachments map[string]string
	// UnfurlLinks attempts to expand a link to show a preview. Success depends
	// on the webpage having the right markdown.
	UnfurlLinks bool
	// UnfurlMedia attempts to expand a link to show a preview. Success depends
	// on the webpage having the right markdown.
	UnfurlMedia bool
	// Username given to bot. If AsUser is true, then message will try to be
	// sent from the given user.
	Username string
	// AsUser attempt to send a message as the user in Username.
	AsUser bool
	// IconURL is a URL to set as the user icon.
	IconURL string
	// IconEmoji is an emoji to set as the user icon.
	IconEmoji string

	Client *http.Client
}

// Send triggers a Slack notification.
func (n *Notification) Send() error {
	if n.AppURL == "" {
		if n.Token == "" {
			return errors.New("missing authentication token or App URL")
		}
		if n.Channel == "" {
			return errors.New("missing channel, group, or username destination")
		}
	}
	if n.Text == "" {
		return errors.New("missing message text")
	}

	attach, err := json.Marshal(n.Attachments)
	if err != nil {
		return err
	}

	// Legacy token-based integration.
	if n.AppURL == "" {
		vals := make(url.Values)
		vals.Set("token", n.Token)
		vals.Set("channel", n.Channel)
		vals.Set("text", n.Text)
		vals.Set("parse", n.Parse)
		vals.Set("link_names", fmt.Sprint(n.LinkNames))
		vals.Set("attachments", string(attach))
		vals.Set("unfurl_links", fmt.Sprintf("%t", n.UnfurlLinks))
		vals.Set("unfurl_media", fmt.Sprintf("%t", n.UnfurlMedia))
		vals.Set("username", n.Username)
		vals.Set("icon_url", n.IconURL)
		vals.Set("icon_emoji", n.IconEmoji)

		resp, err := n.Client.PostForm(API, vals)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var r apiResponse
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			return err
		}

		if !r.OK {
			return errors.New(r.Error)
		}

		return nil
	}

	// New Slack app URL integration.
	data, err := json.Marshal(struct {
		Text string `json:"text"`
	}{n.Text})
	if err != nil {
		return err
	}

	resp, err := n.Client.Post(n.AppURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buff := new(bytes.Buffer)
	_, err = buff.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	s := buff.String()

	if s != "ok" {
		return fmt.Errorf("slack api error: %s", s)
	}

	return nil
}
