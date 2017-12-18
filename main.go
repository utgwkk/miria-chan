package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	miria := NewMiriaClient()
	miria.InitializeSlackClient(os.Getenv("SLACK_WEBHOOK_URL"))
	err := miria.SlackClient.postMessage("Hello from Golang")
	if err != nil {
		panic(err)
	}
}
