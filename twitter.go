package main

import (
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func NewTwitterClient(consumerKey, consumerSecret, accessToken, accessTokenSecret string) *twitter.Client {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	return twitter.NewClient(httpClient)
}

// TweetURL generates tweet's url from IDStr (id_str) and scrrenName (screen_name)
func TweetURL(IDStr, screenName string) string {
	return fmt.Sprintf("https://twitter.com/%s/status/%s", screenName, IDStr)
}
