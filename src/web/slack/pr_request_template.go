package slack

import "encoding/json"

const MessageTemplate = `
{
	"icon_emoji": ":ghost:",
	"username": "Notificador de Prs",
	"blocks": [
		{
			"type": "header",
			"text": {
				"type": "plain_text",
				"text": "Bom dia, aqui vai uma lista de PR's quentinhos para você revisar",
				"emoji": true
			}
		},
		{{ if .High }}
		{
			"type": "divider"
		},
		{
			"type": "context",
			"elements": [
				{
					"type": "mrkdwn",
					"text": ":rotating_light:"
				},
				{
					"type": "mrkdwn",
					"text": "*High Priority*"
				},
				{
					"type": "mrkdwn",
					"text": ":rotating_light:"
				}
			]
		},
		{
			"type": "divider"
		},
		{{ range $index, $element := .High }}
		{{ if $index }},{{ end }}{{ toJson $element }}
		{{ end }},
		{{ end }}
		{{ if .Medium }}
		{
			"type": "divider"
		},
		{
			"type": "context",
			"elements": [
				{
					"type": "image",
					"image_url": "https://api.slack.com/img/blocks/bkb_template_images/highpriority.png",
					"alt_text": "palm tree"
				},
				{
					"type": "mrkdwn",
					"text": "*Medium Priority*"
				}
			]
		},
		{
			"type": "divider"
		},
		{{ range $index, $element := .Medium }}
		{{ if $index }},{{ end }}{{ toJson $element }}
		{{ end }},
		{{ end }}
		{{ if .Low }}
		{
			"type": "divider"
		},
		{
			"type": "context",
			"elements": [
				{
					"type": "image",
					"image_url": "https://api.slack.com/img/blocks/bkb_template_images/mediumpriority.png",
					"alt_text": "palm tree"
				},
				{
					"type": "mrkdwn",
					"text": "*Low Priority (ou sem prioridade)*"
				}
			]
		},
		{
			"type": "divider"
		},
		{{ range $index, $element := .Low }}
		{{ if $index }},{{ end }}{{ toJson $element }}
		{{ end }},
		{{ end }}
		{{ if .ReadyForDeploy }}
		{
			"type": "divider"
		},
		{
			"type": "context",
			"elements": [
				{
					"type": "mrkdwn",
					"text": ":large_green_circle:"
				},
				{
					"type": "mrkdwn",
					"text": "*ReadyForDeploy*"
				}
			]
		},
		{
			"type": "divider"
		},
		{{ range $index, $element := .ReadyForDeploy }}
		{{ if $index }},{{ end }}{{ toJson $element }}
		{{ end }}
		{{ end }}
	]
}
`

func toJson(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}
