package config

import (
	"io/ioutil"
	"label-it/internal/common"
	"log"

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

// YamlConfigV1 interface used to unmarshal YAML configuration
type YamlConfigV1 struct {
	APIVersion int `yaml:"apiVersion"`
	Access     struct {
		User  string `yaml:"user"`
		Token string `yaml:"token"`
	} `yaml:"access"`
	Owner string              `yaml:"owner"`
	Repo  string              `yaml:"repo"`
	Rules map[string]YamlRule `yaml:"rules"`
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

	parseerr := yaml.Unmarshal(dat, &YamlConfig)
	common.CheckErr(parseerr)

	validateVersion(YamlConfig.APIVersion)
}
