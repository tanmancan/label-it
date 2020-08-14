package gitapi

import (
	"encoding/json"
)

// PrBranch properties describing pull request head and base branch
type PrBranch struct {
	Label string `json:"label"`
	Ref   string `json:"ref"`
	SHA   string `json:"sha"`
}

// PrUser properties describing user who opened the pull request
type PrUser struct {
	Login string `json:"login"`
}

// PullRequest individual pull request properties
type PullRequest struct {
	URL    string   `json:"url"`
	Number int      `json:"number"`
	State  string   `json:"state"`
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Head   PrBranch `json:"head"`
	Base   PrBranch `json:"base"`
	Labels []struct {
		Name string `json:"name"`
	} `json:"labels"`
	User PrUser `json:"user"`
}

// ListPullsResponse interface used to unmarshal JSON response
type ListPullsResponse []PullRequest

// ListPulls get a list of open pull requests
// https://docs.github.com/en/rest/reference/pulls#list-pull-requests
func ListPulls() ListPullsResponse {
	endpoint := buildEndpoint(githubConfig.Endpoints.ListPulls)

	query := map[string]string{
		"state": "open",
	}

	request := buildRequest("GET", endpoint, nil, query)
	parsedResponse := gitClient(request)

	prList := ListPullsResponse{}
	json.Unmarshal(parsedResponse, &prList)

	return prList
}
