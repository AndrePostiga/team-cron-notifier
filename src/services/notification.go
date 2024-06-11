package services

import (
	"context"
	"github.com/andrepostiga/team-cron-notifier/src/domain/pullRequest"
	"github.com/andrepostiga/team-cron-notifier/src/domain/team"
)

type InputData struct {
	Teams []struct {
		Name                string `yaml:"name"`
		NotificationConfigs struct {
			Slack struct {
				WebhookSecretEnvName string `yaml:"webhook_secret_env_name"`
			} `yaml:"slack"`
		} `yaml:"notification_configs"`
		Features struct {
			PrNotification struct {
				GithubTokenEnvName string   `yaml:"github_token_env_name"`
				Urls               []string `yaml:"urls"`
			} `yaml:"pr_notification"`
			HealthCheckNotification struct {
				Urls []string `yaml:"urls"`
			} `yaml:"health_check_notification"`
		} `yaml:"features"`
	} `yaml:"teams"`
}

type NotificationService struct {
	PrService PrService
}

type PrService interface {
	GetPrsToNotify(ctx context.Context, team team.Team) ([]pullRequest.PullRequest, error)
	NotifyPrs(ctx context.Context, pullRequests []pullRequest.PullRequest, team team.Team) error
}

func NewNotificationService(prService PrService) *NotificationService {
	return &NotificationService{
		PrService: prService,
	}
}

func (job *NotificationService) Start(ctx context.Context, data *InputData) error {

	for _, teamData := range data.Teams {
		myTeam := team.NewTeam(
			teamData.Name,
			teamData.Features.PrNotification.Urls,
			teamData.Features.HealthCheckNotification.Urls,
			team.NewGeneralSettings(team.NewSlackConfig(teamData.NotificationConfigs.Slack.WebhookSecretEnvName)),
			teamData.Features.PrNotification.GithubTokenEnvName,
		)

		prs, err := job.PrService.GetPrsToNotify(ctx, myTeam)
		if err != nil {
			return err
		}

		err = job.PrService.NotifyPrs(ctx, prs, myTeam)
		if err != nil {
			return err
		}
	}

	return nil
}
