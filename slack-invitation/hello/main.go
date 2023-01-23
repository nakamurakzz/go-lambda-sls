package main

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// 型定義
type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context) (Response, error) {
	var buf bytes.Buffer

	testRes := struct {
		Name string
		Age  int
	}{
		Name: "test",
		Age:  20,
	}
	// goの型 → json：Marshal
	// json → goの型：Unmarshal
	body, err := json.Marshal(testRes)
	
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
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

