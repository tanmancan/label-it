package labeler

import (
	"fmt"
	"label-it/internal/common"
	"label-it/internal/config"
	"label-it/internal/gitapi"
	"regexp"
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

// Checks if beginning or end of a string contains an asterisk(*) rune
func isWildCard(s string) bool {
	if s[len(s)-1] == 42 {
		return true
	}
	if s[0] == 42 {
		return true
	}
	return false
}

// MatchHeadRules determines if provided pull request head branch matche the HeadRule
func (r Rule) MatchHeadRules(pr gitapi.PullRequest) bool {
	wildcardRule := isWildCard(r.HeadRules)

	if wildcardRule == true {
		match, err := regexp.Match(r.HeadRules, []byte(pr.Head.Ref))
		common.CheckErr(err)
		return match
	}

	if r.HeadRules == pr.Head.Ref {
		return true
	}

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

	for _, pr := range prList {
		fmt.Println(pr)
		for _, r := range labelRules {
			matchHead := r.MatchHeadRules(pr)
			fmt.Println(matchHead)
		}
	}

	fmt.Println(labelRules)
}
