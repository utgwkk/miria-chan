package main

import (
    "os"
    "github.com/bluele/slack"
)

func main () {
    slackToken := os.Getenv("MIR_SLACK_TOKEN")
    channelName := os.Getenv("MIR_CHANNEL_NAME")
    api := slack.New(slackToken)
    channel, err := api.FindChannelByName(channelName)
    if err != nil {
        panic(err)
    }
    err = api.ChatPostMessage(channel.Id, "hello", nil)
    if err != nil {
        panic(err)
    }
}