package app

import (
	"context"
	"github.com/andrepostiga/team-cron-notifier/src/domain"
	"github.com/andrepostiga/team-cron-notifier/src/services"
	"log"
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
		log.Fatalf("Failed to load config: %v", err)
	}

	yamlFile, err := os.ReadFile("notifications.yaml")
	if err != nil {
		log.Fatalf("error reading config file: %w", err)
	}

	var prNotificationData services.InputData
	err = yaml.Unmarshal(yamlFile, &prNotificationData)
	if err != nil {
		log.Fatalf("error unmarshalling config file: %w", err)
	}

	log.Printf("Config loaded: %+v", cfg)

	// Create External Services
	githubService, err := github.NewGithubService(web.NewClient(cfg.GithubApi.HTTPClient), cfg.GithubApi)
	if err != nil {
		log.Fatalf("Failed to create Github API: %v", err)
	}

	slackService, err := slack.NewSlackService(web.NewClient(cfg.SlackApi.HTTPClient), &cfg.SlackApi)
	if err != nil {
		log.Fatalf("Failed to create Slack API: %v", err)
	}

	// Create Domain Services
	prService := domain.NewPRService(githubService, slackService)

	notificationService := services.NewNotificationService(prService)
	err = notificationService.Start(ctx, &prNotificationData)
	if err != nil {
		log.Fatalf("Job Failed: %v", err)
	}

}
