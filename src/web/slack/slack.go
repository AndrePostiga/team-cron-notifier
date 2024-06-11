package slack

import (
	"bytes"
	"context"
	"fmt"
	"github.com/andrepostiga/team-cron-notifier/src/domain/pullRequest"
	"github.com/andrepostiga/team-cron-notifier/src/domain/team"
	"github.com/andrepostiga/team-cron-notifier/src/seedwork"
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

	items := make(map[string][]interface{}, len(pullRequests))

	for _, pr := range pullRequests {
		items[pr.Priority().ToString()] = append(items[pr.Priority().ToString()], CreateItem(&pr)...)
	}

	tmpl, err := template.New("message").Funcs(template.FuncMap{"toJson": toJson}).Parse(MessageTemplate)
	if err != nil {
		return fmt.Errorf("error parsing slack request template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, items); err != nil {
		return fmt.Errorf("error executing slack request template: %w", err)
	}

	seedwork.PrintIndentedLog(buf)

	slackConfig, err := team.NotificationSettings().GetSlackConfig()
	if err != nil {
		return fmt.Errorf("Failed to retrieve slack configuration from team: %v", err)
	}
	slackUrl := slack.baseUrl.String() + "/" + os.Getenv(slackConfig.GetWebhookSecretEnvName())

	req, err := http.NewRequest("POST", slackUrl, &buf)
	if err != nil {
		return fmt.Errorf("Failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := slack.client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to send slack api HTTP request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error on slack api Status code:%d error: %v", resp.StatusCode, err)
	}

	return nil
}
