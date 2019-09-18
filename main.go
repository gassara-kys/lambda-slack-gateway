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

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// parse request
	req, err := newRequestForm(&event)
	if err != nil {
		return serverError(err)
	}
	if err := req.valdate(); err != nil {
		return clientError(http.StatusBadRequest)
	}

	// authorization
	ok, err := authorize(req)
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
		res, err = getAlertCount(req)
		if err != nil {
			return serverError(err)
		}
	case "delete":
		print("delete")
		res, err = deleteAlert(req)
		if err != nil {
			return serverError(err)
		}
	default:
		res = getUsage()
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
