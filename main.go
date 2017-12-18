package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	slackClient := NewSlackWebhookClient(os.Getenv("SLACK_WEBHOOK_URL"))
	err := slackClient.postMessage("Hello from Golang")
	if err != nil {
		panic(err)
	}
}
