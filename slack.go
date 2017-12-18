package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type SlackPayload struct {
	Text     string `json:"text"`
	Username string `json:"username"`
}

type SlackWebhookClient struct {
	WebhookURL string
}

func NewSlackWebhookClient(webhookURL string) *SlackWebhookClient {
	return &SlackWebhookClient{webhookURL}
}

func (client *SlackWebhookClient) postMessage(text string) error {
	payload, err := json.Marshal(SlackPayload{
		text,
		"gopher",
	})
	if err != nil {
		return err
	}

	resp, err := http.PostForm(
		client.WebhookURL,
		url.Values{"payload": {string(payload)}},
	)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println(string(body))

	return err
}
