package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/kelseyhightower/envconfig"
)

type dynamoConfig struct {
	Region       string `default:"ap-northeast-1"`               // AWS_DYNAMO_REGION
	AlertTable   string `default:"sns_alert" split_words:"true"` // AWS_DYNAMO_ALERT_TABLE
	AlertHashkey string `default:"timestamp" split_words:"true"` // AWS_DYNAMO_ALERT_HASHKEY
}

type alertTable struct {
	Timestamp string `dynamo:"timestamp"`
	Event     string `dynamo:"event"`
	Message   string `dynamo:"message"`
}

func getAlertTable() (*dynamo.Table, error) {
	var conf dynamoConfig
	if err := envconfig.Process("AWS_DYNAMO", &conf); err != nil {
		return &dynamo.Table{}, err
	}
	table := setupDB(conf.Region, conf.AlertTable)
	return table, nil
}

func setupDB(region, tableNm string) *dynamo.Table {
	db := dynamo.New(session.New(), &aws.Config{Region: aws.String(region)})
	table := db.Table(tableNm)
	return &table
}

func getAlertHashKey() (string, error) {
	var conf dynamoConfig
	if err := envconfig.Process("AWS_DYNAMO", &conf); err != nil {
		return "", err
	}
	return conf.AlertHashkey, nil
}
