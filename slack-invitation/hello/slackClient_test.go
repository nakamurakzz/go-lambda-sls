package main

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func Test_postMessage(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Error("Error loading .env file")
	}

	channelToken := os.Getenv("SLACK_TOKEN")

	slackClient := SlackClient{
		channelToken: channelToken,
		channelName: os.Getenv("SLACK_CHANNEL_NAME"),
		botUserName: "Uwwwwwwww",
	}

	isBotUser := slackClient.isBotUser("Uww")
	if isBotUser {
		t.Error("isBotUser:", isBotUser)
	}
}

func Test_parseSlackChallengeRequest(t *testing.T) {
	reqBody := `{"token": "xxxx", "challenge": "xxxx", "type": "url_verification"}`
	req, err := parseSlackChallengeRequest(reqBody)
	if err != nil {
		t.Error(err)
	}
	if req.Token != "xxxx" || req.Challenge != "xxxx" || req.Type != "url_verification"{
		t.Error("req:", req)
	}
	if !req.isChallenge() {
		t.Error("req:", req)
	}
}

func Test_parseSlackMessageRequest(t *testing.T) {
	reqBody := `{"token": "xxxx", "event" : { "text": "yyyy", "user": "Uzzzzzzzz" } , "type": "event_callback"}`
	req, err := parseSlackChallengeRequest(reqBody)
	if err != nil {
		t.Error(err)
	}
	if req.Event.User != "Uzzzzzzzz" || req.Event.Text != "yyyy" || req.Type != "event_callback"{
		t.Error("req:", req)
	}
	if req.isChallenge() {
		t.Error("req:", req)
	}
}

func Test_getChannelList(t *testing.T) {
	godotenv.Load("../.env")
	channelToken := os.Getenv("SLACK_TOKEN")
	slackClient := SlackClient{
		channelToken: channelToken,
		channelName: os.Getenv("SLACK_CHANNEL_NAME"),
		botUserName: "Uwwwwwwww",
	}
	slackClient.inviteToChannel("xxxxxxxx")
}
