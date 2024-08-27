package github

import (
	"github.com/andrepostiga/team-cron-notifier/src/domain/pullRequest"
	"time"
)

func MapPullRequestsToEntity(response PullRequest) []pullRequest.PullRequest {

	pullRequests := make([]pullRequest.PullRequest, len(response.PullRequestNodes))

	for i, pr := range response.PullRequestNodes {
		createdAt, _ := time.Parse(time.RFC3339, pr.CreatedAt)
		state := pullRequest.State(pr.State)
		author := pullRequest.NewAuthor(pr.Author.Login, pr.Author.AvatarUrl)
		repository := pullRequest.NewRepository(pr.Repository.Name)
		reviews := mapReviewFromGetRepositoryResponse(pr.Reviews)
		labels := mapLabelsFromGetRepositoryResponse(pr.Labels)

		pullRequests[i] = pullRequest.NewPullRequest(
			pr.Title,
			uint(pr.Number),
			pr.URL,
			author,
			createdAt,
			state,
			pr.IsDraft,
			labels,
			reviews,
			repository,
		)
	}

	return pullRequests

}

func mapReviewFromGetRepositoryResponse(response Reviews) []pullRequest.Review {

	reviews := make([]pullRequest.Review, len(response.Nodes))
	for i, review := range response.Nodes {

		if review.State == "COMMENTED" {
			continue
		}

		reviews[i] = pullRequest.NewReview(pullRequest.ReviewState(review.State))
	}

	return reviews
}

func mapLabelsFromGetRepositoryResponse(response Labels) []string {

	labels := make([]string, len(response.Nodes))
	for i, label := range response.Nodes {
		labels[i] = label.Name
	}

	return labels
}
