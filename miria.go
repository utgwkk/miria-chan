package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
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
	if event.Event != "favorite" {
		return
	}
	// If you favorited a tweet and not saved the images yet
	tweetID := event.TargetObject.IDStr
	tweetUser := event.TargetObject.User.ScreenName
	tweetURL := TweetURL(tweetID, tweetUser)
	log.Printf("You favorited %s.", tweetURL)
	if !m.shouldBeSaved(event) {
		return
	}
	err := m.SlackClient.postMessage(tweetURL)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *MiriaClient) PostYourFavoritedTweetWithMediaAndSaveImages(event *twitter.Event) {
	if event.Event != "favorite" {
		return
	}
	tweetID := event.TargetObject.IDStr
	tweetUser := event.TargetObject.User.ScreenName
	tweetURL := TweetURL(tweetID, tweetUser)
	log.Printf("You favorited %s.", tweetURL)
	if !m.shouldBeSaved(event) {
		return
	}
	err := m.SlackClient.postMessage(tweetURL)
	if err != nil {
		log.Fatal(err)
	}
	tempDir, err := ioutil.TempDir("", "miria")
	if err != nil {
		log.Fatal(err)
	}
	medias := event.TargetObject.ExtendedEntities.Media
	for _, media := range medias {
		// Save image to temporary directory
		downloadURL := media.MediaURLHttps
		filename := path.Base(downloadURL)
		destinationPath := path.Join(tempDir, filename)
		err := download(media.MediaURLHttps, destinationPath)
		if err != nil {
			log.Print(err)
		}
		// TODO: Save information to DB
		// TODO: Generate a thumbnail
		// TODO: Save image to S3 bucket
	}
}

func (m *MiriaClient) saveInfoToDB(tweet *twitter.Tweet, filename string) {
	res, err := m.DB.Exec("insert into images (filename) values (?)", filename)
	if err != nil {
		log.Fatal(err)
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	tweetID := tweet.IDStr
	tweetUser := tweet.User.ScreenName
	tweetURL := TweetURL(tweetID, tweetUser)
	comment := tweet.FullText
	_, err = m.DB.Exec(
		"insert into image_info (image_id, comment, source) values (?, ?, ?)",
		lastID, comment, tweetURL,
	)
	if err != nil {
		log.Fatal(err)
	}
}

func download(url, destination string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(body)
	return nil
}

func (m *MiriaClient) shouldBeSaved(event *twitter.Event) bool {
	tweetID := event.TargetObject.IDStr
	tweetUser := event.TargetObject.User.ScreenName
	tweetURL := TweetURL(tweetID, tweetUser)
	hasMedia := len(event.TargetObject.ExtendedEntities.Media) > 0
	return event.Source.IDStr == m.TwitterUserID && hasMedia && !m.existSource(tweetURL)
}

func (m *MiriaClient) existSource(source string) bool {
	rows, err := m.DB.Query("select count(*) from image_info where source = ?", source)
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
