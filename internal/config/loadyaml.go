package config

import (
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/tanmancan/label-it/v1/internal/common"

	"gopkg.in/yaml.v2"
)

// YamlConfig loaded yaml configuration
var YamlConfig YamlConfigV1

// RuleTypeString groups of rule types for string values.
// Exact - the rule value must be an exact match of the compare value.
// NoExact - a rule value must NOT be an exact match of the compare value.
// Match - a regex pattern that must match a compare value.
// NoMatch - a regex pattern that must NOT match a compare value
type RuleTypeString struct {
	Exact   string `yaml:"exact,omitempty"`
	NoExact string `yaml:"no-exact,omitempty"`
	Match   string `yaml:"match,omitempty"`
	NoMatch string `yaml:"no-match,omitempty"`
}

// RuleTypeInt groups of rule types for integer values
// Exact - the rule value must be an exact match of the compare value.
// NoExact - a rule value must NOT be an exact match of the compare value.
// Match - a regex pattern that must match a compare value.
// NoMatch - a regex pattern that must NOT match a compare value
type RuleTypeInt struct {
	Exact   int    `yaml:"exact,omitempty"`
	NoExact int    `yaml:"no-exact,omitempty"`
	Match   string `yaml:"match,omitempty"`
	NoMatch string `yaml:"no-match,omitempty"`
}

// YamlRuleGroup rules for an individual label
type YamlRuleGroup struct {
	Label  string         `yaml:"label"`
	Head   RuleTypeString `yaml:"head-rule,omitempty"`
	Base   RuleTypeString `yaml:"base-rule,omitempty"`
	Title  RuleTypeString `yaml:"title-rule,omitempty"`
	Body   RuleTypeString `yaml:"body-rule,omitempty"`
	User   RuleTypeString `yaml:"user-rule,omitempty"`
	Number RuleTypeInt    `yaml:"number-rule,omitempty"`
}

// YamlGithubAccess stores user and access token for Github api authentication
type YamlGithubAccess struct {
	User  string `yaml:"user"`
	Token string `yaml:"token"`
}

// Parses access values and checks for env variables if provided
func parseAccess(v string) (string, error) {
	if v[0] == 36 {
		env := v[1:]
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
	APIVersion int              `yaml:"apiVersion"`
	Access     YamlGithubAccess `yaml:"access"`
	Owner      string           `yaml:"owner"`
	Repo       string           `yaml:"repo"`
	Rules      []YamlRuleGroup  `yaml:"rules"`
}

// Validates YAML with current package version
func validateVersion(ver int) {
	if ver != ConfigVersion {
		log.Fatal("Invalid config file version. Current tool requires version", ConfigVersion)
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
