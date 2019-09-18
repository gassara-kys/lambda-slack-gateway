package main

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
)

type requestForm struct {
	TeamDomain       string `json:"team_domain"`
	ChannelName      string `json:"channel_name"`
	UserName         string `json:"user_name"`
	Command          string `json:"command"`
	Text             string `json:"text"`
	RequestBody      string `json:"request_body"`              // all request body
	RequestTimestamp string `json:"x_slack_request_timestamp"` // X-Slack-Request-Timestamp header
	SlackSignature   string `json:"x_slack_signature"`         // X-Slack-Signature header
}

func newRequestForm(event *events.APIGatewayProxyRequest) (*requestForm, error) {
	var form requestForm
	form.RequestTimestamp = event.Headers["X-Slack-Request-Timestamp"]
	form.SlackSignature = event.Headers["X-Slack-Signature"]
	form.RequestBody = event.Body
	values, err := url.ParseQuery(form.RequestBody)
	if err != nil {
		return &form, fmt.Errorf("parse query error :%#v", err)
	}
	form.TeamDomain = values.Get("team_domain")
	form.ChannelName = values.Get("channel_name")
	form.UserName = values.Get("user_name")
	form.Command = values.Get("command")
	form.Text = values.Get("text")

	return &form, nil
}

func (r *requestForm) valdate() error {
	if r.Command == "" ||
		r.RequestBody == "" ||
		r.RequestTimestamp == "" ||
		r.SlackSignature == "" {
		// `request_body` contains sensitive data, log output is not possible.
		fmt.Printf("[ERROR]The following parameters are required: "+
			"command=%s, request_body=..., x_slack_request_timestamp=%s, x_slack_signature=%s",
			r.Command, r.RequestTimestamp, r.SlackSignature)
		return errors.New("request parameters(command, requestbody, timestamp, signature) are required")
	}
	return nil
}
