package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strconv"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type slackConfig struct {
	SigningSecret string `default:"TEST_KEY" split_words:"true"` // SLACK_SIGNING_SECRET
	SkipTimeValid string `default:"" split_words:"true"`         // SLACK_SKIP_TIME_VALID
}

func authrize(req request) (bool, error) {
	var conf slackConfig
	if err := envconfig.Process("slack", &conf); err != nil {
		return false, err
	}
	// valid request time
	reqUnix, err := strconv.ParseInt(req.RequestTimestamp, 10, 64)
	if err != nil {
		return false, err
	}
	max := time.Now().Unix()
	min := time.Now().Add(-5 * time.Minute).Unix()
	if conf.SkipTimeValid == "true" {
		min = 0
	}
	log.Printf("[INFO]RequestTimeVlidation: min=%d, reqUnix=%d, max=%d", min, reqUnix, max)
	if reqUnix < min || reqUnix > max {
		return false, nil
	}

	// valid signature
	base := "v0:" + req.RequestTimestamp + ":" + req.RequestBody
	signature := "v0=" + makeHMACSha256(base, conf.SigningSecret)
	if signature != req.SlackSignature {
		return false, nil
	}

	log.Printf("[Authorized] user=%s(%s), channel=%s, cmd=%s(%s)",
		req.UserName, req.TeamDomain, req.ChannelName, req.Command, req.Text)
	return true, nil
}

func makeHMACSha256(msg, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(msg))
	return hex.EncodeToString(mac.Sum(nil))
}
