package requestInMemoryImplementation

import (
	"github.com/andrepostiga/team-cron-notifier/src/domain/pullRequest"
	"github.com/andrepostiga/team-cron-notifier/src/web/slack"
)

type SlackApiRequest struct {
	IconEmoji string        `json:"icon_emoji"`
	Username  string        `json:"username"`
	Blocks    []interface{} `json:"blocks"`
}

type Block struct {
	Type string `json:"type"`
}

type Header struct {
	Block
	Text struct {
		Type  string `json:"type"`
		Text  string `json:"text"`
		Emoji bool   `json:"emoji"`
	} `json:"text"`
}

type Divider struct {
	Block
}

type PrioritySection struct {
	Block
	Elements []Block `json:"elements"`
}

type MarkDown struct {
	Block
	Text string `json:"text"`
}

func (req *SlackApiRequest) AddBlock(block map[string]interface{}) {
	req.Blocks = append(req.Blocks, block)
}

func BuildRequest(prs []pullRequest.PullRequest, slackMoji string, messageUserName string) SlackApiRequest {

	items := make(map[pullRequest.Priority][]interface{}, len(prs))
	for _, pr := range prs {
		items[pr.Priority()] = append(items[pr.Priority()], slack.CreateItem(&pr)...)
	}

	req := SlackApiRequest{
		IconEmoji: slackMoji,
		Username:  messageUserName,
		Blocks:    make([]interface{}, 0),
	}

	req.AddBlock(map[string]interface{}{
		"type": "header",
		"text": map[string]interface{}{
			"type":  "plain_text",
			"text":  "Bom dia, aqui vai uma lista de PR's quentinhos para você revisar",
			"emoji": true,
		},
	})

	req.AddBlock(DividerBlock)
	req.AddBlock(map[string]interface{}{
		"type":     "context",
		"elements": HighPriorityBlock,
	})
	req.AddBlock(DividerBlock)

	// Append high priority PRs
	for _, item := range items[pullRequest.High] {
		req.AddBlock(item.(map[string]interface{}))
	}

	req.AddBlock(DividerBlock)
	req.AddBlock(map[string]interface{}{
		"type":     "context",
		"elements": MediumPriorityBlock,
	})
	req.AddBlock(DividerBlock)

	// Prs Medium
	for _, item := range items[pullRequest.Medium] {
		req.AddBlock(item.(map[string]interface{}))
	}

	req.AddBlock(DividerBlock)
	req.AddBlock(map[string]interface{}{
		"type":     "context",
		"elements": LowPriorityBlock,
	})
	req.AddBlock(DividerBlock)

	// Prs Low
	for _, item := range items[pullRequest.Low] {
		req.AddBlock(item.(map[string]interface{}))
	}

	return req
}
