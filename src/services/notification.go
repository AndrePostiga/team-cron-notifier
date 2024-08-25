package services

import (
	"context"
	"github.com/andrepostiga/team-cron-notifier/src/domain/pullRequest"
	"github.com/andrepostiga/team-cron-notifier/src/domain/team"
	"log/slog"
)

type InputData struct {
	Teams []Teams `yaml:"teams"`
}

type Teams struct {
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
}

func (t Teams) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("Name", t.Name),
		slog.Group("NotificationConfigs",
			slog.Group("Slack",
				slog.String("WebhookSecretEnvName", "******"),
			),
		),
		slog.Group("Features",
			slog.Group("PrNotification",
				slog.String("GithubTokenEnvName", "******"),
				slog.Any("Urls", t.Features.PrNotification.Urls),
			),
			slog.Group("HealthCheckNotification",
				slog.Any("Urls", t.Features.HealthCheckNotification.Urls),
			),
		),
	)
}

type NotificationService struct {
	PrService PrService
	logger    *slog.Logger
}

type PrService interface {
	GetPrsToNotify(ctx context.Context, team team.Team) ([]pullRequest.PullRequest, error)
	NotifyPrs(ctx context.Context, pullRequests []pullRequest.PullRequest, team team.Team) error
}

func NewNotificationService(logger *slog.Logger, prService PrService) *NotificationService {
	return &NotificationService{
		PrService: prService,
		logger:    logger,
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

		job.logger.Info("job is starting for team", "teamName", myTeam.Name(), "content", teamData)

		prs, err := job.PrService.GetPrsToNotify(ctx, myTeam)
		if err != nil {
			job.logger.ErrorContext(ctx, "error when try to get prs from vcs", "error", err)
			return err
		}

		err = job.PrService.NotifyPrs(ctx, prs, myTeam)
		if err != nil {
			job.logger.ErrorContext(ctx, "error when try to notify prs", "error", err)
			return err
		}
	}

	return nil
}
