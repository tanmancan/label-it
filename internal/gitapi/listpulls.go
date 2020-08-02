package gitapi

import (
	"encoding/json"
)

// PullRequest individual pull request properties
type PullRequest struct {
	URL    string `json:"url"`
	Number int    `json:"number"`
	State  string `json:"state"`
	Title  string `json:"title"`
	Head   struct {
		Label string `json:"label"`
		Ref   string `json:"ref"`
		SHA   string `json:"sha"`
	} `json:"head"`
	Base struct {
		Label string `json:"label"`
		Ref   string `json:"ref"`
		SHA   string `json:"sha"`
	} `json:"base"`
	Labels []struct {
		Name string `json:"name"`
	} `json:"labels"`
}

// ListPullsResponse interface used to unmarshal JSON response
type ListPullsResponse []PullRequest

// ListPulls get a list of open pull requests
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
