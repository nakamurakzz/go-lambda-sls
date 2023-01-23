package main

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func Test_postMessage(t *testing.T) {
	err := godotenv.Load(".env")
	if err != nil {
		t.Error("Error loading .env file")
	}

	channelToken := os.Getenv("SLACK_TOKEN")
	channelName := os.Getenv("SLACK_CHANNEL_NAME")
	slackClient := SlackClient{
		channelToken: channelToken,
		channelName: channelName,
	}
	res :=slackClient.postMessage("test message from lambda")
	if res != nil {
		t.Error("token:", channelToken)
		t.Error(res)
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
}