package app

import (
	"context"
	"fmt"
	"github.com/andrepostiga/team-cron-notifier/src/domain"
	"github.com/andrepostiga/team-cron-notifier/src/services"
	"log"
	"log/slog"
	"os"

	"github.com/andrepostiga/team-cron-notifier/src/config"
	"github.com/andrepostiga/team-cron-notifier/src/web"
	"github.com/andrepostiga/team-cron-notifier/src/web/github"
	"github.com/andrepostiga/team-cron-notifier/src/web/slack"
	"gopkg.in/yaml.v3"
)

func BuildJob(ctx context.Context) {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("Failed to load config: %v", err))
	}

	logLevel := new(slog.LevelVar)
	err = logLevel.UnmarshalText([]byte(cfg.LogLevel))
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	}))
	slog.SetDefault(logger)

	yamlFile, err := os.ReadFile("notifications.yaml")
	if err != nil {
		logger.Error("error reading config file", "error", err)
		panic(err)
	}

	var prNotificationData services.InputData
	err = yaml.Unmarshal(yamlFile, &prNotificationData)
	if err != nil {
		logger.Error("error unmarshalling config file", "error", err)
		panic(err)
	}

	// Create External Services
	githubService, err := github.NewGithubService(logger, web.NewClient(logger, cfg.GithubApi.HTTPClient), cfg.GithubApi)
	if err != nil {
		logger.Error("failed to create Github API", "error", err)
		panic(err)
	}

	slackService, err := slack.NewSlackService(web.NewClient(logger, cfg.SlackApi.HTTPClient), &cfg.SlackApi)
	if err != nil {
		logger.Error("failed to create Slack API", "error", err)
		panic(err)
	}

	// Create Domain Services
	prService := domain.NewPRService(logger, githubService, slackService)

	notificationService := services.NewNotificationService(logger, prService)
	err = notificationService.Start(ctx, &prNotificationData)
	if err != nil {
		logger.ErrorContext(ctx, "Job Failed", "error", err)
		panic(err)
	}

}
