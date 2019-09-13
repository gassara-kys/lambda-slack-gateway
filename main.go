package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)
var responseHeader = map[string]string{
	"Content-Type": "application/json",
}

type request struct {
	TeamDomain       string `json:"team_domain"`
	ChannelName      string `json:"channel_name"`
	UserName         string `json:"user_name"`
	Command          string `json:"command"`
	Text             string `json:"text"`
	RequestBody      string `json:"request_body"`              // all request body
	RequestTimestamp string `json:"x_slack_request_timestamp"` // X-Slack-Request-Timestamp header
	SlackSignature   string `json:"x_slack_signature"`         // X-Slack-Signature header
}

func handler(ctx context.Context, req request) (events.APIGatewayProxyResponse, error) {
	log.Printf("START lambda function")
	// authorization
	ok, err := authrize(req)
	if err != nil {
		return serverError(err)
	}
	if !ok {
		return clientError(http.StatusForbidden)
	}

	// handle request param
	var res *map[string]interface{}
	switch req.Text {
	case "get":
		res, err = getAlertCount()
		if err != nil {
			return serverError(err)
		}
	case "delete":
		print("delete")
		res, err = deleteAlert()
		if err != nil {
			return serverError(err)
		}
	default:
		usage := map[string]interface{}{
			"usage": "hogehogehogehoge",
		}
		res = &usage
	}

	resJSON, err := json.Marshal(res)
	if err != nil {
		return serverError(err)
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    responseHeader,
		Body:       string(resJSON),
	}, nil
}

// 5xx error
func serverError(err error) (events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

// 4xx error
func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func main() {
	lambda.Start(handler)
}
