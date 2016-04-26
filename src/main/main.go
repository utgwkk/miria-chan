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
    options := &slack.ChatPostMessageOpt{
        Username: "赤城みりあ",
        IconUrl: "https://gyazo.com/04a9d82be786486aa2c44463e9e5b60d.png"}
    err = api.ChatPostMessage(channel.Id, "みりあもやるー！", options)
    if err != nil {
        panic(err)
    }
}