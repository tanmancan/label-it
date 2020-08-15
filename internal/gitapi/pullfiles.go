package gitapi

import (
	"encoding/json"
	"sort"
	"strconv"
	"sync"
)

// PrFile properties describing a changed file in a pull request
type PrFile struct {
	SHA      string `json:"sha"`
	Filename string `json:"filename"`
	Status   string `json:"status"`
}

// ListPrFilesResponse A list of files from the list pull request endpoint
type ListPrFilesResponse []PrFile

// Get a list of changed files for a given pull request number.
// Github only shows a maximum of 100 files per page, so we are
// recursively calling the api page by page until we get to the last
// page with files
// https://docs.github.com/en/rest/reference/pulls#list-pull-requests-files
func listPrFiles(number int, page int, wg *sync.WaitGroup, c chan ListPrFilesResponse) {
	endpoint := buildEndpoint(githubConfig.Endpoints.ListPrFiles, number)

	perPage := 100
	query := map[string]string{
		"per_page": strconv.Itoa(perPage),
		"page":     strconv.Itoa(page),
	}

	request := buildRequest("GET", endpoint, nil, query)

	parsedResponse := gitClient(request)

	prFiles := ListPrFilesResponse{}
	json.Unmarshal(parsedResponse, &prFiles)
	c <- prFiles

	// Hardcoding 500 max files (100 per page)
	// Possibly make this configarable, but need to find out
	// what the constraints are on the list pr files endpoint
	maxFiles := 1000
	maxPages := maxFiles / perPage
	if len(prFiles) == 100 && page <= maxPages-1 {
		page++
		wg.Add(1)
		go listPrFiles(number, page, wg, c)
	}
	wg.Done()
}

// GetAllFiles calls the get pr files endpoint recursively for
// each page that return a list of files
func GetAllFiles(number int) []string {

	c := make(chan ListPrFilesResponse)
	var wg sync.WaitGroup
	wg.Add(1)
	go listPrFiles(number, 1, &wg, c)
	go func() {
		wg.Wait()
		close(c)
	}()

	var allFiles ListPrFilesResponse

	for res := range c {
		allFiles = append(allFiles, res...)
	}

	var allFileNames []string
	for _, files := range allFiles {
		allFileNames = append(allFileNames, files.Filename)
	}
	sort.Strings(allFileNames)
	return allFileNames
}
