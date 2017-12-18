package main

import "github.com/dghubble/go-twitter/twitter"

type MiriaClient struct {
	TwitterClient (*twitter.Client)
	SlackClient   (*SlackWebhookClient)
}

func NewMiriaClient() *MiriaClient {
	return &MiriaClient{}
}

func (m *MiriaClient) InitializeTwitterClient(consumerKey, consumerSecret, accessToken, accessTokenSecret string) {
	m.TwitterClient = NewTwitterClient(consumerKey, consumerSecret, accessToken, accessTokenSecret)
}

func (m *MiriaClient) InitializeSlackClient(webhookURL string) {
	m.SlackClient = NewSlackWebhookClient(webhookURL)
}
