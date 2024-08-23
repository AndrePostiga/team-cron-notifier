package domain

import (
	"context"
	"fmt"
	"github.com/andrepostiga/team-cron-notifier/src/domain/pullRequest"
	"github.com/andrepostiga/team-cron-notifier/src/domain/team"
	"log/slog"
)

type PullRequestFetcher interface {
	GetPullRequests(ctx context.Context, team team.Team) ([]pullRequest.PullRequest, error)
}

type NotifyService interface {
	SendMessage(ctx context.Context, pullRequests []pullRequest.PullRequest, team team.Team) error
}

type PRService struct {
	logger                    *slog.Logger
	pullRequestFetcherService PullRequestFetcher
	notifyService             NotifyService
}

func NewPRService(logger *slog.Logger, pullRequestFetcherService PullRequestFetcher, notifyService NotifyService) *PRService {
	return &PRService{
		pullRequestFetcherService: pullRequestFetcherService,
		notifyService:             notifyService,
		logger:                    logger,
	}
}

func (svc *PRService) GetPrsToNotify(ctx context.Context, team team.Team) ([]pullRequest.PullRequest, error) {
	if !team.IsPrNotificationEnabled() {
		return nil, fmt.Errorf("PR notification is disabled for team %s", team.Name())
	}

	pullRequests, err := svc.pullRequestFetcherService.GetPullRequests(ctx, team)
	if err != nil {
		return nil, err
	}

	return pullRequests, nil
}

func (svc *PRService) NotifyPrs(ctx context.Context, pullRequests []pullRequest.PullRequest, team team.Team) error {
	if !team.IsPrNotificationEnabled() {
		return fmt.Errorf("PR notification is disabled for team %s", team.Name())
	}

	err := svc.notifyService.SendMessage(ctx, pullRequests, team)
	if err != nil {
		return err
	}

	return nil
}
