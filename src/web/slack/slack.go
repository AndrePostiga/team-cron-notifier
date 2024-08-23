package slack

import (
	"bytes"
	"context"
	"fmt"
	"github.com/andrepostiga/team-cron-notifier/src/domain/pullRequest"
	"github.com/andrepostiga/team-cron-notifier/src/domain/team"
	"net/http"
	"net/url"
	"os"
	"text/template"

	"github.com/andrepostiga/team-cron-notifier/src/config"
)

type Slack struct {
	client  *http.Client
	baseUrl *url.URL
}

func NewSlackService(client *http.Client, ops *config.SlackApiConfig) (*Slack, error) {
	baseUrl, err := url.Parse(ops.BaseUrl)
	if err != nil {
		return nil, err
	}

	return &Slack{
		client:  client,
		baseUrl: baseUrl,
	}, nil
}

func (slack *Slack) SendMessage(ctx context.Context, pullRequests []pullRequest.PullRequest, team team.Team) error {

	slackConfig, err := team.NotificationSettings().GetSlackConfig()
	if err != nil {
		return fmt.Errorf("failed to retrieve slack configuration from team: %w", err)
	}

	message, err := createMessage(pullRequests)
	if err != nil {
		return fmt.Errorf("failed to create slack message from template: %w", err)
	}

	err = slack.doRequest(ctx, slackConfig, message)
	if err != nil {
		return fmt.Errorf("failed to send message to slack: %w", err)
	}

	return nil
}

func createMessage(pullRequests []pullRequest.PullRequest) (*bytes.Buffer, error) {
	items := make(map[string][]interface{}, len(pullRequests))

	for _, pr := range pullRequests {
		items[pr.Priority().ToString()] = append(items[pr.Priority().ToString()], CreateItem(&pr)...)
	}

	tmpl, err := template.New("message").Funcs(template.FuncMap{"toJson": toJson}).Parse(MessageTemplate)
	if err != nil {
		return nil, fmt.Errorf("error parsing slack request template: %w", err)
	}

	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, items); err != nil {
		return nil, fmt.Errorf("error executing slack request template: %w", err)
	}

	return buf, nil
}

func (slack *Slack) doRequest(ctx context.Context, slackConfig team.SlackConfig, message *bytes.Buffer) error {
	slackUrl := slack.baseUrl.String() + "/" + os.Getenv(slackConfig.GetWebhookSecretEnvName())

	req, err := http.NewRequestWithContext(ctx, "POST", slackUrl, message)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := slack.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send slack api HTTP request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error on slack api Status code:%d error: %v", resp.StatusCode, err)
	}

	return nil
}
