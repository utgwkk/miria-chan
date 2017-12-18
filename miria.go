package main

import (
	"database/sql"
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
	DB            (*sql.DB)
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

func (m *MiriaClient) InitializeDBConnection(hostname, databaseName, username, password string) {
	db, err := NewMySQLConnection(hostname, databaseName, username, password)
	if err != nil {
		log.Fatal(err)
	}
	m.DB = db
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

func (m *MiriaClient) JustPostYourFavoritedTweetWithMediaWhenNotSavedYet(event *twitter.Event) {
	// If you favorited a tweet and not saved the images yet
	if m.shouldBeSaved(event) {
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

func (m *MiriaClient) shouldBeSaved(event *twitter.Event) bool {
	tweetID := event.TargetObject.IDStr
	tweetUser := event.TargetObject.User.ScreenName
	tweetURL := TweetURL(tweetID, tweetUser)
	hasMedia := len(event.TargetObject.ExtendedEntities.Media) > 0
	return event.Event == "favorite" && event.Source.IDStr == m.TwitterUserID && hasMedia && !m.existSource(tweetURL)
}

func (m *MiriaClient) existSource(source string) bool {
	rows, err := m.DB.Query("select count(*) from image_info where source = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var count int64
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Fatal(err)
		}
	}
	return count > 0
}
