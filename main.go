package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	miria := NewMiriaClient()
	miria.RegisterThumbnailPath(os.Getenv("THUMBNAIL_DIR"))
	miria.InitializeSlackClient(os.Getenv("SLACK_WEBHOOK_URL"))
	miria.SlackClient.SetUsername(os.Getenv("SLACK_USERNAME"))
	miria.SlackClient.SetIconEmoji(os.Getenv("SLACK_ICON_EMOJI"))
	miria.InitializeTwitterClient(
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"),
		os.Getenv("TWITTER_ACCESS_TOKEN_KEY"),
		os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
	)
	miria.InitializeAWSCredential(
		os.Getenv("AWS_ACCESS_KEY_ID"),
		os.Getenv("AWS_SECRET_ACCESS_KEY"),
		os.Getenv("AWS_REGION"),
		os.Getenv("AWS_BUCKET"),
		os.Getenv("AWS_BUCKET_BASEPATH"),
	)
	miria.InitializeDBConnection(
		os.Getenv("DB_HOSTNAME"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
	)
	miria.CollectEvents(miria.PostYourFavoritedTweetWithMediaAndSaveImages)
}
