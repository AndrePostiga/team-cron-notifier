package github

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

const (
	// Github GraphQl Query
	GetPullRequestsRequest string = `
	query (
		$owner: String!
		$repo: String!
		$reviewStates: [PullRequestReviewState!]
	) {
		repository(owner: $owner, name: $repo) {
			pullRequests(first: 50, states: OPEN) {
				nodes {
					title
					number
					url
					author {
						login
						avatarUrl
					}
					createdAt
					state
					reviews(first: 50, states: $reviewStates) {
						nodes {
							state
						}
					}
					labels(first: 10) {
						nodes {
							name
						}
					}
					repository {
						name
					}
				}
			}
		}
	}
	`
)
