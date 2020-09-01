package twilio

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type apiResponse struct {
	Sid                 string      `json:"sid"`
	DateCreated         string      `json:"date_created"`
	DateUpdated         string      `json:"date_updated"`
	DateSent            interface{} `json:"date_sent"`
	AccountSid          string      `json:"account_sid"`
	To                  string      `json:"to"`
	From                string      `json:"from"`
	MessagingServiceSid interface{} `json:"messaging_service_sid"`
	Body                string      `json:"body"`
	Status              string      `json:"status"`
	NumSegments         string      `json:"num_segments"`
	NumMedia            string      `json:"num_media"`
	Direction           string      `json:"direction"`
	APIVersion          string      `json:"api_version"`
	Price               interface{} `json:"price"`
	PriceUnit           string      `json:"price_unit"`
	ErrorCode           string      `json:"error_code"`
	ErrorMessage        string      `json:"error_message"`
	URI                 string      `json:"uri"`
	SubresourceUris     struct {
		Media string `json:"media"`
	} `json:"subresource_uris"`
}

// Notification is Twilio notification
type Notification struct {
	Content    string `json:"content"`
	NumberTo   string `json:"numberto"`
	NumberFrom string `json:"numberfrom"`
	AccountSid string `json:"accountsid"`
	AuthToken  string `json:"authtoken"`
}

// API is the Twilio API endpoint.
var API = "https://api.twilio.com"

// Send sends a Twilio notification.
func (n *Notification) Send() error {
	if n.NumberTo == "" {
		return errors.New("missing receiver number")
	}

	if n.NumberFrom == "" {
		return errors.New("missing sender number")
	}

	if n.AccountSid == "" {
		return errors.New("missing account sid")
	}

	if n.AuthToken == "" {
		return errors.New("missing auth token")
	}

	msgData := url.Values{}
	msgData.Set("To", n.NumberTo)
	msgData.Set("From", n.NumberFrom)
	msgData.Set("Body", n.Content)
	msgDataReader := *strings.NewReader(msgData.Encode())

	url := fmt.Sprintf("%s/2010-04-01/Accounts/%s/Messages.json", API, n.AccountSid)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, &msgDataReader)
	req.SetBasicAuth(n.AccountSid, n.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	if r.ErrorCode != "" {
		return errors.New(r.ErrorMessage)
	}

	return nil
}
