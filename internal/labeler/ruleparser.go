package labeler

import (
	"fmt"
	"label-it/internal/config"
	"label-it/internal/gitapi"
)

// Rule rule from YAML config
type Rule struct {
	Label      string
	HeadRules  string
	BaseRules  string
	TitleRules []string
}

// LabelRules set of rules created from YAML config
type LabelRules []Rule

// MatchHeadRules determines if provided pull request head branch matche the HeadRule
func (r Rule) MatchHeadRules(pr gitapi.PullRequest) bool {
	return false
}

// MatchBaseRules determines if provided pull request base branch matche theBaseRule
func (r Rule) MatchBaseRules(pr gitapi.PullRequest) bool {
	return false
}

// MatchTitleRules determines if provided pull request contains text in title rules
func (r Rule) MatchTitleRules(pr gitapi.PullRequest) bool {
	return false
}

// RuleParser parses YAML config to create LabelRules
func RuleParser(prList gitapi.ListPullsResponse) {
	labelRules := LabelRules{}

	// fmt.Println(labelRules)
	for _, rule := range config.YamlConfig.Rules {
		newRule := Rule{rule.Label, rule.Head, rule.Base, rule.Title}
		labelRules = append(labelRules, newRule)
	}

	fmt.Println(labelRules)
}
