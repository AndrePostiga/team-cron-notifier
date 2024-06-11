package slack

import (
	"fmt"
	"github.com/andrepostiga/team-cron-notifier/src/domain/pullRequest"
)

func CreateItem(pr *pullRequest.PullRequest) []interface{} {
	section := map[string]interface{}{
		"type": "section",
		"text": map[string]interface{}{
			"type": "mrkdwn",
			"text": fmt.Sprintf("*<%s|%s>*", pr.URL(), pr.Title()),
		},
		"accessory": map[string]interface{}{
			"type": "button",
			"text": map[string]interface{}{
				"type":  "plain_text",
				"text":  fmt.Sprintf("%s", pr.Repository().Name()),
				"emoji": true,
			},
			"value": "click_me_123",
			"url":   pr.URL(),
		},
	}

	author := pr.Author()

	context := map[string]interface{}{
		"type": "context",
		"elements": []interface{}{
			map[string]interface{}{"type": "mrkdwn", "text": fmt.Sprintf(":white_check_mark::white_check_mark: %d Approves", pr.NumberOfApproves())},
			map[string]interface{}{"type": "mrkdwn", "text": fmt.Sprintf(":no_entry: %d RequestChanges", pr.NumberOfRequestChanges())},
			map[string]interface{}{"type": "mrkdwn", "text": fmt.Sprintf("Aberto há *%d dias*", pr.GetOpenedDays())},
			map[string]interface{}{"type": "image", "image_url": author.AvatarUrl(), "alt_text": author.Name()},
			map[string]interface{}{"type": "mrkdwn", "text": fmt.Sprintf("%s", author.Name())},
		},
	}

	return []interface{}{section, context}
}
