package main

import (
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

var event = events.APIGatewayProxyRequest{
	Headers: map[string]string{
		"X-Slack-Request-Timestamp": "1500000000",
		"X-Slack-Signature":         "xxxxxxxxxx",
	},
	Body: "team_domain=team.com&channel_name=ch_test&user_name=alice&command=/test&text=param",
}

var form = requestForm{
	TeamDomain:       "team.com",
	ChannelName:      "ch_test",
	UserName:         "alice",
	Command:          "/test",
	Text:             "param",
	RequestBody:      "team_domain=team.com&channel_name=ch_test&user_name=alice&command=/test&text=param",
	RequestTimestamp: "1500000000",
	SlackSignature:   "xxxxxxxxxx",
}

func TestNewRequestForm(t *testing.T) {
	tests := []struct {
		name  string
		input *events.APIGatewayProxyRequest
		want  *requestForm
	}{
		{"Basic", &event, &form},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := newRequestForm(test.input)
			if err != nil {
				t.Fatalf("want no err, but has error %#v", err)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("want %#v, but %#v", test.want, got)
			}
		})
	}
}

func TestValdate(t *testing.T) {
	tests := []struct {
		name  string
		input *requestForm
		want  error
	}{
		{"Basic", &form, nil},
		{"RequiredError", &requestForm{}, errors.New("request parameters(command, requestbody, timestamp, signature) are required")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.valdate()
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("want %#v, but %#v", test.want, got)
			}
		})
	}
}
