package main

import (
	"github.com/nlopes/slack"
)

func getUsage() *map[string]interface{} {
	return &map[string]interface{}{
		"response_type": "ephemeral",
		"attachments": []slack.Attachment{
			slack.Attachment{
				Title: "usage",
				Color: "#216546",
				Fields: []slack.AttachmentField{
					slack.AttachmentField{
						Title: "/alert get",
						Value: "return number of alert recourds.",
					},
					slack.AttachmentField{
						Title: "/alert delete",
						Value: "clear alert records",
					},
				},
			},
		},
	}
}
