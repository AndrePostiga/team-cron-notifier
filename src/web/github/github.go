package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/andrepostiga/team-cron-notifier/src/config"
	"github.com/andrepostiga/team-cron-notifier/src/domain/pullRequest"
	"github.com/andrepostiga/team-cron-notifier/src/domain/team"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Github struct {
	client  *http.Client
	baseUrl *url.URL
}

func NewGithubService(client *http.Client, ops config.GithubApiConfig) (*Github, error) {
	baseUrl, err := url.Parse(ops.BaseUrl)
	if err != nil {
		return nil, err
	}

	return &Github{
		client:  client,
		baseUrl: baseUrl,
	}, nil
}

// Rate limit of this endpoint using GRAPH_QL_Api:
// Primary: Individual calls cannot request more than 500,000 total nodes.
// How many Nodes are returned in this query GetPullRequestsRequest?
// 50                = 50 pull requests
// 50 x 50		     = 2500 reviews
// 50 x 50 x 10      = 25000 labels
// 25000 + 2500 + 50 = 27550 nodes is under 500.000 limit
// The query has 127 points
// Each gpg_token can query 5000 points per hour
// So, we can query 5000/127 = 39 times per hour (a.k.a 39 repositories)
// Secondary:
// https://docs.github.com/en/graphql/overview/rate-limits-and-node-limits-for-the-graphql-api#rate-limit
func (git *Github) GetPullRequests(ctx context.Context, team team.Team) ([]pullRequest.PullRequest, error) {
	prs := make([]pullRequest.PullRequest, 0)

	for _, repo := range team.PrNotification().Repositories() {
		prsInRepo, err := git.doRequest(ctx, repo, team)
		if err != nil {
			return nil, fmt.Errorf("failed to get pull requests for %s from github: %w", repo, err)
		}

		prs = append(prs, prsInRepo...)
	}

	return prs, nil
}

func (git *Github) doRequest(ctx context.Context, orgRepositoryName string, team team.Team) ([]pullRequest.PullRequest, error) {
	parts := strings.Split(orgRepositoryName, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid repository name, it should be in the format 'org/repository' requested:%s", orgRepositoryName)
	}

	requestPayload := GraphQLRequest{
		Query: GetPullRequestsRequest,
		Variables: map[string]interface{}{
			"owner":        parts[0],
			"repo":         parts[1],
			"reviewStates": []string{"APPROVED", "CHANGES_REQUESTED"},
		},
	}

	requestBody, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal github graphqlrequest payload: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", git.baseUrl.String(), bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request_in_memory_implementation: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv(team.GithubApiToken()))

	resp, err := git.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request_in_memory_implementation: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var respData GetPullRequestsResponse
	if err := json.Unmarshal(body, &respData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	prs := MapPullRequestsToEntity(respData.Data.Repository.PullRequests)
	return prs, nil
}
