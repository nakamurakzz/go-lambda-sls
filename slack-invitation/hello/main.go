package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
	"github.com/slack-go/slack"
)

var (
	secretCache, _ = secretcache.New()
)

// 型定義
type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {

	// slackからのリクエストを構造体にマッピング
	var slackRequest SlackRequest
	slackRequest, parseError := parseSlackChallengeRequest(request.Body)
	if parseError != nil {
		log.Println(parseError)
		return Response{StatusCode: 404}, parseError
	}

	// challengeリクエストの場合はchallengeを返す
	if slackRequest.isChallenge() {
		body, marshalError := json.Marshal(struct {
			Challenge string
		}{
			Challenge: slackRequest.Challenge,
		})
		if marshalError != nil {
			log.Println(marshalError)
			return Response{StatusCode: 404}, marshalError
		}

		return Response{
			StatusCode:      200,
			IsBase64Encoded: false,
			Body:            string(body),
			Headers: map[string]string{
				"Content-Type":           "application/json",
				"X-MyCompany-Func-Reply": "hello-handler",
			},
		}, nil
	}

	// シークレットを取得
	slackToken, secretCacheError := secretCache.GetSecretString(os.Getenv("SSM_KEY_NAME"))
	if secretCacheError != nil {
		log.Println(secretCacheError)
		return Response{StatusCode: 404}, secretCacheError
	}

	
	// Slackにメッセージを投稿
	slackClient := SlackClient{
		channelToken: slackToken,
		channelName: os.Getenv("SLACK_CHANNEL_NAME"),
		botUserName: os.Getenv("SLACK_BOT_USER_NAME"),
	}

	if slackClient.isBotUser(slackRequest.Event.User) {
		log.Println("bot user")
		return Response{StatusCode: 200}, nil
	}

	// Slackにメッセージを投稿
	slackErr := slackClient.postMessage("test message from lambda")
	if slackErr != nil {
		log.Println(slackErr)
		return Response{StatusCode: 500}, slackErr
	}

	// goの型 → json：Marshal
	// json → goの型：Unmarshal
	body, toJsonErr := json.Marshal(struct {
		message string
	}{
		message: "ok",
	})
	
	if toJsonErr != nil {
		// errをログ表示
		log.Println(toJsonErr)
		return Response{StatusCode: 404}, toJsonErr
	}

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(body),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}

	return resp, nil
}

func main() {
	// ハンドラーを引数として渡して実行する
	lambda.Start(Handler)
}

type Event struct {
	Text string `json:"text"`
	User string `json:"user"`
}

type SlackRequest struct {
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	Type      string `json:"type"`
	Event		 Event  `json:"event"`
}

func (c SlackRequest) isChallenge() bool {
	return c.Type == "url_verification"
}

func parseSlackChallengeRequest(body string) (SlackRequest, error) {
	var request SlackRequest
	err := json.Unmarshal([]byte(body), &request)
	log.Println(request)
	return request, err
}

type SlackClient struct {
	channelToken string
	channelName	string
	botUserName	string
}

func (c SlackClient) postMessage (message string) error {
	log.Println("Post Slack Message")
	client := slack.New(c.channelToken)
	log.Println("post to channel: ", c.channelName)
	_, _, err := client.PostMessage(c.channelName, slack.MsgOptionText(message, false))
	return err
}

func (c SlackClient) isBotUser(userName string) bool {
	if c.botUserName == "" || userName == "" {
		return false
	}
	res := userName == c.botUserName
	log.Println("is Bot User?",res, userName, c.botUserName)
	return userName == c.botUserName
}