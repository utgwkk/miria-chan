package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type SlackPayload struct {
	Text     string `json:"text"`
	Username string `json:"username"`
}

func postMessage(text string) error {
	payload, err := json.Marshal(SlackPayload{
		text,
		"gopher",
	})
	if err != nil {
		return err
	}

	webhookUrl := os.Getenv("SLACK_WEBHOOK_URL")

	resp, err := http.PostForm(
		webhookUrl,
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
