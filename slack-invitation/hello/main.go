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
	var slackChallengeRequest SlackChallengeRequest
	slackChallengeRequest, parseError := parseSlackChallengeRequest(request.Body)
	if parseError != nil {
		log.Println(parseError)
		return Response{StatusCode: 404}, parseError
	}

	if slackChallengeRequest.isChallenge() {
		// challengeの場合はchallengeを返す
		body, marshalError := json.Marshal(struct {
			Challenge string
		}{
			Challenge: slackChallengeRequest.Challenge,
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
		// チャンネル名を取得
		slackChannelName, slackChannelError := secretCache.GetSecretString(os.Getenv("SSM_KEY_NAME"))
		if slackChannelError != nil {
			log.Println(secretCacheError)
			return Response{StatusCode: 404}, secretCacheError
		}
		
	
	slackClient := SlackClient{
		channelToken: slackToken,
		channelName: slackChannelName,
	}

	slackErr := slackClient.postMessage("test message from lambda")

	if slackErr != nil {
		log.Println(slackErr)
		return Response{StatusCode: 500}, slackErr
	}

	// var buf bytes.Buffer
	res := struct {
		message string
	}{
		message: "ok",
	}
	// goの型 → json：Marshal
	// json → goの型：Unmarshal
	body, toJsonErr := json.Marshal(res)
	
	if toJsonErr != nil {
		// errをログ表示
		log.Println(toJsonErr)
		return Response{StatusCode: 404}, toJsonErr
	}

	// json.HTMLEscape(&buf, body)
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

type SlackChallengeRequest struct {
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	Type      string `json:"type"`
}

func (c SlackChallengeRequest) isChallenge() bool {
	return c.Type == "url_verification"
}

func parseSlackChallengeRequest(body string) (SlackChallengeRequest, error) {
	var request SlackChallengeRequest
	err := json.Unmarshal([]byte(body), &request)
	log. Println(request)
	return request, err
}

type SlackClient struct {
	channelToken string
	channelName	string
}

func (c SlackClient) postMessage (message string) error {
	log.Println("Post Slack Message")
	client := slack.New(c.channelToken)
	log.Println(c.channelName, client)
	// _, _, err := client.PostMessage(c.channelName, slack.MsgOptionText(message, false))
	return nil
}