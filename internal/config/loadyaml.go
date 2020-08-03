package config

import (
	"errors"
	"io/ioutil"
	"label-it/internal/common"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// YamlConfig loaded yaml configuration
var YamlConfig YamlConfigV1

// YamlRule an individual rule for a label
type YamlRule struct {
	Head   string `yaml:"head,omitempty"`
	Base   string `yaml:"base,omitempty"`
	Title  string `yaml:"title,omitempty"`
	Body   string `yaml:"body,omitempty"`
	User   string `yaml:"user,omitempty"`
	Number []int  `yaml:"number,omitempty"`
}

// YamlGithubAccess stores user and access token for Github api authentication
type YamlGithubAccess struct {
	User  string `yaml:"user"`
	Token string `yaml:"token"`
}

// Parses access values and checks for env variables if provided
func parseAccess(v string) (string, error) {
	if v[0] == 36 {
		env := v[1:len(v)]
		envval := os.Getenv(env)

		if envval == "" {
			return "", errors.New("Env variable not found")
		}

		return envval, nil
	}

	return v, nil
}

// UnmarshalYAML custom parser for access data in YAML
func (a *YamlGithubAccess) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var access struct {
		User  string
		Token string
	}

	err := unmarshal(&access)

	if err != nil {
		return err
	}

	if access.User == "" {
		return errors.New("Missing access user")
	}

	if access.Token == "" {
		return errors.New("Missing access token")
	}

	parsedUser, userErr := parseAccess(access.User)

	if userErr != nil {
		return userErr
	}

	parsedToken, tokenErr := parseAccess(access.Token)

	if tokenErr != nil {
		return tokenErr
	}

	a.User = parsedUser
	a.Token = parsedToken
	return nil
}

// YamlConfigV1 interface used to unmarshal YAML configuration
type YamlConfigV1 struct {
	APIVersion int                 `yaml:"apiVersion"`
	Access     YamlGithubAccess    `yaml:"access"`
	Owner      string              `yaml:"owner"`
	Repo       string              `yaml:"repo"`
	Rules      map[string]YamlRule `yaml:"rules"`
}

// Validates YAML with current package version
func validateVersion(ver int) {
	if ver != AppConfig.ConfigVersion {
		log.Fatal("Invalid config file version. Current tool requires version", AppConfig.ConfigVersion)
	}
}

// LoadYaml load configuration from a given yaml file
func LoadYaml() {
	dat, err := ioutil.ReadFile(YamlPath)
	common.CheckErr(err)

	parseerr := yaml.UnmarshalStrict(dat, &YamlConfig)
	common.CheckErr(parseerr)

	validateVersion(YamlConfig.APIVersion)
}
