package gitapi

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"

	"github.com/tanmancan/label-it/v1/internal/config"
)

func Test_buildEndpoint(t *testing.T) {
	*&config.YamlConfig.Owner = "world"
	*&config.YamlConfig.Repo = "Robot"
	type args struct {
		endpointTemplate string
		args             []interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"replace strings and digits",
			args{
				"Hello %[1]s, my name is %[2]s #%[3]d. Why did the %[4]s cross the network?",
				[]interface{}{
					42,
					"packet",
				},
			},
			"Hello world, my name is Robot #42. Why did the packet cross the network?",
		},
		{
			"creates add labels endpoint",
			args{
				githubConfig.Endpoints.AddLabels,
				[]interface{}{
					44,
				},
			},
			"/repos/world/Robot/issues/44/labels",
		},
		{
			"creates list pulls endpoint",
			args{
				githubConfig.Endpoints.ListPulls,
				[]interface{}{},
			},
			"/repos/world/Robot/pulls",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildEndpoint(tt.args.endpointTemplate, tt.args.args...); got != tt.want {
				t.Errorf("buildEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildAPIURL(t *testing.T) {
	type args struct {
		endpoint string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"build full add label url",
			args{
				"/repos/Hello/World/issues/55/labels",
			},
			"https://api.github.com/repos/Hello/World/issues/55/labels",
		},
		{
			"build fake api url",
			args{
				"/abc/def/hij",
			},
			"https://api.github.com/abc/def/hij",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildAPIURL(tt.args.endpoint); got != tt.want {
				t.Errorf("buildAPIURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildBasicAuth(t *testing.T) {
	fakeuser := "fakeuser"
	faketoken := "faketoken"
	*&config.YamlConfig.Access.User = fakeuser
	*&config.YamlConfig.Access.Token = faketoken
	authtoken := fmt.Sprintf("%[1]s:%[2]s", fakeuser, faketoken)
	fakebasicauth := base64.StdEncoding.EncodeToString([]byte(authtoken))
	fakebasicauth = fmt.Sprintf("Basic %[1]s", fakebasicauth)
	tests := []struct {
		name string
		want string
	}{
		{
			"correctly encodes username and token",
			fakebasicauth,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildBasicAuth(); got != tt.want {
				t.Errorf("buildBasicAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildReqQuery(t *testing.T) {
	request, _ := http.NewRequest("GET", "http://example.com", nil)
	query := map[string]string{
		"testq":    "testval",
		"anotherq": "secondval",
	}
	want := "http://example.com?anotherq=secondval&testq=testval"
	t.Run("query params added to request", func(t *testing.T) {
		buildReqQuery(request, query)
		reqURL := request.URL.String()
		if reqURL != want {
			t.Errorf("request url = %v, want %v", reqURL, want)
		}
		fmt.Println(reqURL)
	})
}
