package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type SlackPayload struct {
	Text      string `json:"text"`
	Username  string `json:"username"`
	IconURL   string `json:"icon_url"`
	IconEmoji string `json:"icon_emoji"`
}

type SlackWebhookClient struct {
	WebhookURL string
	Username   string
	IconURL    string
	IconEmoji  string
}

func NewSlackWebhookClient(webhookURL string) *SlackWebhookClient {
	return &SlackWebhookClient{webhookURL, "", "", ""}
}

func (sc *SlackWebhookClient) SetUsername(username string) {
	sc.Username = username
}

func (sc *SlackWebhookClient) SetIconURL(iconURL string) {
	sc.IconURL = iconURL
}

func (sc *SlackWebhookClient) SetIconEmoji(iconEmoji string) {
	sc.IconEmoji = iconEmoji
}

func (client *SlackWebhookClient) postMessage(text string) error {
	payload, err := json.Marshal(SlackPayload{
		text,
		client.Username,
		client.IconURL,
		client.IconEmoji,
	})
	if err != nil {
		return err
	}
	log.Println(string(payload))

	resp, err := http.PostForm(
		client.WebhookURL,
		url.Values{"payload": {string(payload)}},
	)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	bodyStr := string(body)
	if bodyStr != "ok" {
		return errors.New(bodyStr)
	}

	return nil
}
