package github

type GetPullRequestsResponse struct {
	Data struct {
		Repository Repository `json:"repository"`
	} `json:"data"`
	Errors []struct {
		Type    string   `json:"type"`
		Path    []string `json:"path"`
		Message string   `json:"message"`
	} `json:"errors"`
}

type Repository struct {
	PullRequests PullRequest `json:"pullRequests"`
}

type PullRequest struct {
	PullRequestNodes []PullRequestNodes `json:"nodes"`
}

type PullRequestNodes struct {
	Title  string `json:"title"`
	Number int    `json:"number"`
	URL    string `json:"url"`
	Author struct {
		Login     string `json:"login"`
		AvatarUrl string `json:"avatarUrl"`
	} `json:"author"`
	CreatedAt  string  `json:"createdAt"`
	State      string  `json:"state"`
	IsDraft    bool    `json:"isDraft"`
	Reviews    Reviews `json:"reviews"`
	Labels     Labels  `json:"labels"`
	Repository struct {
		Name string `json:"name"`
	} `json:"repository"`
}

type Reviews struct {
	Nodes []struct {
		State string `json:"state"`
	} `json:"nodes"`
}

type Labels struct {
	Nodes []struct {
		Name string `json:"name"`
	} `json:"nodes"`
}
