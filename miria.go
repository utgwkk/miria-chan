package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dghubble/go-twitter/twitter"
)

type MiriaClient struct {
	TwitterClient (*twitter.Client)
	TwitterUserID string
	SlackClient   (*SlackWebhookClient)
	AWS           (*AWSCredential)
}

func NewMiriaClient() *MiriaClient {
	return &MiriaClient{}
}

func (m *MiriaClient) InitializeTwitterClient(consumerKey, consumerSecret, accessToken, accessTokenSecret string) {
	m.TwitterClient = NewTwitterClient(consumerKey, consumerSecret, accessToken, accessTokenSecret)
	// Get authenticated user's id_str
	user, _, err := m.TwitterClient.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{})
	if err != nil {
		log.Fatal(err)
	}
	m.TwitterUserID = user.IDStr
}

func (m *MiriaClient) InitializeSlackClient(webhookURL string) {
	m.SlackClient = NewSlackWebhookClient(webhookURL)
}

func (m *MiriaClient) InitializeAWSCredential(accessKeyID, secretAccessKey, region, bucketName, basePath string) {
	m.AWS = NewAWSCredential(accessKeyID, secretAccessKey, region, bucketName, basePath)
}

func (m *MiriaClient) CollectEvents(processEvent func(*twitter.Event)) {
	demux := twitter.NewSwitchDemux()
	demux.Event = processEvent
	stream, err := m.TwitterClient.Streams.User(&twitter.StreamUserParams{})
	if err != nil {
		log.Fatal(err)
	}
	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	stream.Stop()
}

// JustPostYourFavoritedTweetToSlack just post your favorited tweet's URL to Slack.
func (m *MiriaClient) JustPostYourFavoritedTweetToSlack(event *twitter.Event) {
	eventKind := event.Event
	eventSource := event.Source.IDStr
	// If you favorited a tweet
	if eventKind == "favorite" && eventSource == m.TwitterUserID {
		tweetID := event.TargetObject.IDStr
		tweetUser := event.TargetObject.User.ScreenName
		tweetURL := TweetURL(tweetID, tweetUser)
		log.Printf("You favorited %s.", tweetURL)
		err := m.SlackClient.postMessage(tweetURL)
		if err != nil {
			log.Fatal(err)
		}
	}
}
