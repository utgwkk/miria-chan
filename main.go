package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	miria := NewMiriaClient()
	miria.InitializeSlackClient(os.Getenv("SLACK_WEBHOOK_URL"))
	miria.SlackClient.SetUsername(os.Getenv("SLACK_USERNAME"))
	miria.SlackClient.SetIconEmoji(os.Getenv("SLACK_ICON_EMOJI"))
	miria.InitializeTwitterClient(
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"),
		os.Getenv("TWITTER_ACCESS_TOKEN_KEY"),
		os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
	)
	miria.CollectEvents(miria.JustPostYourFavoritedTweetToSlack)
}
