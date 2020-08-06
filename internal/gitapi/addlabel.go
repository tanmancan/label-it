package gitapi

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tanmancan/label-it/v1/internal/common"
)

// PrLabel interface describing a pull request and
// a list of labels to add to the pull request
type PrLabel struct {
	Issue  int
	Labels []string
}

// AddLabels adds given list of labels to a specific pull request
// https://docs.github.com/en/rest/reference/issues#set-labels-for-an-issue
func AddLabels(prLabel PrLabel, c chan string) {
	endpoint := buildEndpoint(githubConfig.Endpoints.AddLabels, prLabel.Issue)

	reqBody, err := json.Marshal(map[string][]string{
		"labels": prLabel.Labels,
	})
	common.CheckErr(err)

	request := buildRequest("POST", endpoint, reqBody, nil)
	gitClient(request)

	log := fmt.Sprintf(
		"Added labels \"%[1]s\" to PR #%[2]d",
		strings.Join(prLabel.Labels, ", "),
		prLabel.Issue,
	)
	c <- log
}
