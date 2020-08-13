package gitapi

import (
	"encoding/json"
	"strconv"
)

// PrFile properties describing a changed file in a pull request
type PrFile struct {
	SHA      string `json:"sha"`
	Filename string `json:"filename"`
	Status   string `json:"status"`
}

// ListPrFilesResponse A list of files from the list pull request endpoint
type ListPrFilesResponse []PrFile

// ListPrFiles get a list of changed files for a given pull request number
// https://docs.github.com/en/rest/reference/pulls#list-pull-requests-files
func ListPrFiles(number int, page int) ListPrFilesResponse {
	endpoint := buildEndpoint(githubConfig.Endpoints.ListPrFiles, number)

	query := map[string]string{
		"per_page": "100",
		"page":     strconv.Itoa(page),
	}

	request := buildRequest("GET", endpoint, nil, query)
	parsedResponse := gitClient(request)

	prFiles := ListPrFilesResponse{}
	json.Unmarshal(parsedResponse, &prFiles)

	return prFiles
}
