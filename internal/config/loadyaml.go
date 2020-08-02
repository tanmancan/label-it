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
	Label string   `yaml:"label"`
	Head  string   `yaml:"head,omitempty"`
	Base  string   `yaml:"base,omitempty"`
	Title []string `yaml:"title,omitempty"`
}

// YamlConfigV1 configuration type for YAML unmarshalling
type YamlConfigV1 struct {
	APIVersion int `yaml:"apiVersion"`
	Access     struct {
		User  string `yaml:"user"`
		Token string `yaml:"token"`
	} `yaml:"access"`
	Owner string     `yaml:"owner"`
	Repo  string     `yaml:"repo"`
	Rules []YamlRule `yaml:"rules"`
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
