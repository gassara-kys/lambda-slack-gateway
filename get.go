package main

import (
	"strconv"

	"github.com/nlopes/slack"
)

func getAlertCount(req *requestForm) (*map[string]interface{}, error) {
	var result map[string]interface{}
	table, err := getAlertTable()
	if err != nil {
		return &result, err
	}

	var results []alertTable
	if err := table.Scan().All(&results); err != nil {
		return &result, err
	}
	result = map[string]interface{}{
		"response_type": "in_channel",
		"attachments": []slack.Attachment{
			slack.Attachment{
				Title: "/alert command called",
				Color: "#764FA5",
				Fields: []slack.AttachmentField{
					slack.AttachmentField{
						Title: "operation",
						Value: "get",
					},
					slack.AttachmentField{
						Title: "called by",
						Value: req.UserName,
					},
					slack.AttachmentField{
						Title: "alert_count",
						Value: strconv.Itoa(len(results)),
					},
				},
			},
		},
	}
	return &result, nil
}
