package main

import (
	"log"
	"strconv"

	"github.com/nlopes/slack"
)

func deleteAlert(req *requestForm) (*map[string]interface{}, error) {
	var result map[string]interface{}
	table, err := getAlertTable()
	if err != nil {
		return &result, err
	}
	hashKey, err := getAlertHashKey()
	if err != nil {
		return &result, err
	}

	var records []alertTable
	if err := table.Scan().All(&records); err != nil {
		return &result, err
	}
	for _, record := range records {
		if err := table.Delete(hashKey, record.Timestamp).Run(); err != nil {
			return &result, err
		}
	}
	log.Printf("[DELETE]%d alert deleted.", len(records))
	result = map[string]interface{}{
		"response_type": "in_channel",
		"attachments": []slack.Attachment{
			slack.Attachment{
				Title: "/alert command called",
				Color: "#764FA5",
				Fields: []slack.AttachmentField{
					slack.AttachmentField{
						Title: "operation",
						Value: "delete",
					},
					slack.AttachmentField{
						Title: "called by",
						Value: req.UserName,
					},
					slack.AttachmentField{
						Title: "delete_count",
						Value: strconv.Itoa(len(records)),
					},
				},
			},
		},
	}
	return &result, nil
}
