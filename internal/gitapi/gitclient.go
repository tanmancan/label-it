package gitapi

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"label-it/internal/common"
	"label-it/internal/config"
	"net/http"
)

// Endpoints use by this package
type githubAPIEndpoints struct {
	AddLabels string
	ListPulls string
}

// Configuration types for Github API
type githubAPIConfig struct {
	Version        string
	BaseURL        string
	RequestHeaders map[string]string
	Endpoints      githubAPIEndpoints
}

var githubConfig = githubAPIConfig{
	Version: "v3",
	BaseURL: "https://api.github.com",
	RequestHeaders: map[string]string{
		"Accept":       "application/vnd.github.v3+json",
		"Content-Type": "application/json",
	},
	Endpoints: githubAPIEndpoints{
		AddLabels: "/repos/%[1]s/%[2]s/issues/%[3]d/labels",
		ListPulls: "/repos/%[1]s/%[2]s/pulls",
	},
}

// Indent and prints a JSON response
func prettyPrintResponse(content []byte) {
	dst := &bytes.Buffer{}
	indenterr := json.Indent(dst, content, "", "    ")
	common.CheckErr(indenterr)
	fmt.Println(dst.String())
}

// Populate endpoint templates in AppConfig.Github.Endpoints with provided arguments
func buildEndpoint(endpointTemplate string, args ...interface{}) string {
	owner := config.YamlConfig.Owner
	repo := config.YamlConfig.Repo
	argsWithRepo := []interface{}{owner, repo}

	if len(args) > 0 {
		for _, arg := range args {
			argsWithRepo = append(argsWithRepo, arg)
		}
	}

	return fmt.Sprintf(endpointTemplate, argsWithRepo...)
}

// Generate Github API URL for a given endpoint string
func buildAPIURL(endpoint string) string {
	return fmt.Sprintf("%[1]s%[2]s", githubConfig.BaseURL, endpoint)
}

// Generate Basic Authentication token for a request
func buildBasicAuth() string {
	token := fmt.Sprintf("%[1]s:%[2]s", config.YamlConfig.Access.User, config.YamlConfig.Access.Token)
	tokenenc := base64.StdEncoding.EncodeToString([]byte(token))
	return fmt.Sprintf("Basic %s", tokenenc)
}

// Generate request body
func buildReqBody(reqBody []byte) *bytes.Buffer {
	return bytes.NewBuffer(reqBody)
}

// Generate query string for a request and given parameter
func buildReqQuery(request *http.Request, reqQueryParam map[string]string) {
	if reqQueryParam != nil {
		q := request.URL.Query()
		for key, val := range reqQueryParam {
			q.Add(key, val)
		}
		request.URL.RawQuery = q.Encode()
	}
}

// Builds a API request to be used in http.Client
func buildRequest(method string, endpoint string, reqBody []byte, reqQueryParam map[string]string) *http.Request {
	if method == "" {
		method = "GET"
	}

	if reqBody != nil {
		method = "POST"
	}

	url := buildAPIURL(endpoint)
	body := buildReqBody(reqBody)

	request, reqerr := http.NewRequest(method, url, body)
	common.CheckErr(reqerr)

	for key, value := range githubConfig.RequestHeaders {
		request.Header.Add(key, value)
	}

	authtoken := buildBasicAuth()
	request.Header.Add("Authorization", authtoken)

	buildReqQuery(request, reqQueryParam)

	return request
}

// Client for making http request to Github API
func gitClient(request *http.Request) []byte {
	client := http.Client{}

	res, resperr := client.Do(request)
	common.CheckErr(resperr)

	content, readerr := ioutil.ReadAll(res.Body)
	common.CheckErr(readerr)

	res.Body.Close()

	return content
}
