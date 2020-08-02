package gitapi

import (
	"encoding/json"
	"label-it/internal/common"
)

// AddLabels adds given list of labels to a specific pull request
// https://docs.github.com/en/rest/reference/issues#set-labels-for-an-issue
func AddLabels(issue int, labels []string) {
	endpoint := buildEndpoint(githubConfig.Endpoints.AddLabels, issue)

	reqBody, err := json.Marshal(map[string][]string{
		"labels": labels,
	})
	common.CheckErr(err)

	request := buildRequest("POST", endpoint, reqBody, nil)
	gitClient(request)
}
