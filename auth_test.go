package main

import (
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

var base = requestForm{
	TeamDomain:       "team.com",
	ChannelName:      "ch_test",
	UserName:         "alice",
	Command:          "/test",
	Text:             "param",
	RequestBody:      "team_domain=team.com&channel_name=ch_test&user_name=alice&command=/test&text=param",
	RequestTimestamp: strconv.FormatInt(time.Now().Unix(), 10),
	SlackSignature:   "",
}

func TestAuthorize(t *testing.T) {
	// set env
	os.Setenv("SLACK_SIGNING_SECRET", "TEST_KEY")
	os.Setenv("SLACK_SKIP_TIME_VALID", "false")

	// set signature
	validSignature := "v0=" + makeHMACSha256(
		"v0:"+base.RequestTimestamp+":"+base.RequestBody, os.Getenv("SLACK_SIGNING_SECRET"))
	base.SlackSignature = validSignature

	// error test data
	parseErr := base
	parseErr.RequestTimestamp = "xxx"
	old := base
	old.RequestTimestamp = strconv.FormatInt(time.Now().Add(-5*time.Minute).Unix(), 10)
	future := base
	future.RequestTimestamp = strconv.FormatInt(time.Now().Add(10*time.Minute).Unix(), 10)
	sigErr := base
	sigErr.SlackSignature = "xxxxxxxxxxx"

	// tests
	tests := []struct {
		name      string
		input     *requestForm
		want      bool
		wantError bool
		err       error
	}{
		{"basic", &base, true, false, nil},
		{"parseErr", &parseErr, false, true, &strconv.NumError{Func: "ParseInt"}},
		{"too old", &old, false, false, nil},
		{"future", &future, false, false, nil},
		{"invlalid signature", &sigErr, false, false, nil},
	}

	for _, test := range tests {
		got, err := authorize(test.input)
		if !test.wantError && err != nil {
			t.Fatalf("want no err, but has error %#v", err)
		}
		if test.wantError && reflect.TypeOf(err) != reflect.TypeOf(test.err) {
			t.Fatalf("want %#v, but %#v", test.err, err)
		}
		if got != test.want {
			t.Fatalf("want %t, but %t", test.want, got)
		}
	}
}

func TestMakeHMACSha256(t *testing.T) {
	key := "TEST_KEY"
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"basic", "v0:1500000000:xxx", "850133b9a887eecfd67737407aebca593b4034b1f7156bb9169d8c4fff416477"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := makeHMACSha256(test.input, key)
			if got != test.want {
				t.Fatalf("want: %s, but got: %s", test.want, got)
			}
		})
	}
}
