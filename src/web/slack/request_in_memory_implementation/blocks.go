package requestInMemoryImplementation

var HighPriorityBlock = []interface{}{
	map[string]interface{}{
		"type": "mrkdwn",
		"text": ":rotating_light:",
	},
	map[string]interface{}{
		"type": "mrkdwn",
		"text": "*High Priority*",
	},
	map[string]interface{}{
		"type": "mrkdwn",
		"text": ":rotating_light:",
	},
}

var MediumPriorityBlock = []interface{}{
	map[string]interface{}{
		"type":      "image",
		"image_url": "https://api.slack.com/img/blocks/bkb_template_images/highpriority.png",
		"alt_text":  "medium priority",
	},
	map[string]interface{}{
		"type": "mrkdwn",
		"text": "*Medium Priority*",
	},
}

var LowPriorityBlock = []interface{}{
	map[string]interface{}{
		"type":      "image",
		"image_url": "https://api.slack.com/img/blocks/bkb_template_images/mediumpriority.png",
		"alt_text":  "low priority",
	},
	map[string]interface{}{
		"type": "mrkdwn",
		"text": "*Low Priority*",
	},
}

var DividerBlock = map[string]interface{}{
	"type": "divider",
}
