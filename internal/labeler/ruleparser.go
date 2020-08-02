package labeler

import (
	"label-it/internal/common"
	"label-it/internal/config"
	"label-it/internal/gitapi"
	"regexp"
	"sort"
)

// Rule rule from YAML config
type Rule struct {
	Label       string
	HeadRules   string
	BaseRules   string
	TitleRules  string
	BodyRules   string
	UserRule    string
	NumberRules []int
}

// LabelRules set of rules created from YAML config
type LabelRules []Rule

// MatchHeadRules determines if provided pull request head branch matche the HeadRule
func (r Rule) MatchHeadRules(pr gitapi.PullRequest) bool {
	if r.HeadRules == pr.Head.Ref {
		return true
	}

	// Fall back to regular expression search if no exact match found for rule and branch name
	// Should return true if no rule provide
	// TODO: validate claim above
	match, err := regexp.MatchString(r.HeadRules, pr.Head.Ref)
	common.CheckErr(err)

	return match
}

// MatchBaseRules determines if provided pull request base branch matche theBaseRule
func (r Rule) MatchBaseRules(pr gitapi.PullRequest) bool {
	if r.BaseRules == pr.Base.Ref {
		return true
	}

	// Fall back to regular expression search if no exact match found for rule and branch name
	// Should return true if no rule provide
	// TODO: validate claim above
	match, err := regexp.MatchString(r.BaseRules, pr.Base.Ref)
	common.CheckErr(err)
	return match
}

// MatchTitleRules determines if provided pull request contains text in title rules
func (r Rule) MatchTitleRules(pr gitapi.PullRequest) bool {
	if r.TitleRules == "" {
		return true
	}

	match, err := regexp.MatchString(r.TitleRules, pr.Title)
	common.CheckErr(err)

	return match
}

// MatchBodyRules determines if provided pull request contains text in title rules
func (r Rule) MatchBodyRules(pr gitapi.PullRequest) bool {
	if r.BodyRules == "" {
		return true
	}

	match, err := regexp.MatchString(r.BodyRules, pr.Body)
	common.CheckErr(err)

	return match
}

// MatchUserRules checks if pull request creator username matches user rule
func (r Rule) MatchUserRules(pr gitapi.PullRequest) bool {
	if r.UserRule == "" {
		return true
	}

	return r.UserRule == pr.User.Login
}

// MatchNumberRules determines if pull request issue number matches provider number in rule
func (r Rule) MatchNumberRules(pr gitapi.PullRequest) bool {
	if len(r.NumberRules) == 0 {
		return true
	}

	searchIdx := sort.SearchInts(r.NumberRules, pr.Number)

	if searchIdx == len(r.NumberRules) {
		return false
	}

	return r.NumberRules[searchIdx] == pr.Number
}

// MatchAllRules checks if a pull request passes all checks for a given rule
func (r Rule) MatchAllRules(pr gitapi.PullRequest) bool {
	switch {
	case r.MatchHeadRules(pr) != true:
		return false
	case r.MatchBaseRules(pr) != true:
		return false
	case r.MatchTitleRules(pr) != true:
		return false
	case r.MatchBodyRules(pr) != true:
		return false
	case r.MatchUserRules(pr) != true:
		return false
	case r.MatchNumberRules(pr) != true:
		return false
	}

	return true
}

// Checks if pull request already has label
func prHasLabel(pr gitapi.PullRequest, label string) bool {
	if len(pr.Labels) == 0 {
		return false
	}

	var existingLabels []string

	for _, prlabel := range pr.Labels {
		existingLabels = append(existingLabels, prlabel.Name)
	}

	sort.Strings(existingLabels)
	searchIdx := sort.SearchStrings(existingLabels, label)

	if searchIdx == len(existingLabels) {
		return false
	}

	return existingLabels[searchIdx] == label
}

// RuleParser parses rules and checks if they match provided pull requests
// returns a list of matched pull request numbers and labels to apply to them
func RuleParser(prList gitapi.ListPullsResponse) []gitapi.PrLabel {
	labelRules := LabelRules{}

	for label, rule := range config.YamlConfig.Rules {
		numRule := rule.Number
		sort.Ints(numRule)
		newRule := Rule{label, rule.Head, rule.Base, rule.Title, rule.Body, rule.User, numRule}
		labelRules = append(labelRules, newRule)
	}

	matchedLabelPr := []gitapi.PrLabel{}

	for _, pr := range prList {
		newLabels := []string{}
		for _, r := range labelRules {
			hasLabel := prHasLabel(pr, r.Label)

			if hasLabel == false {
				matchAll := r.MatchAllRules(pr)
				if matchAll == true {
					newLabels = append(newLabels, r.Label)
				}
			}
		}

		if len(newLabels) != 0 {
			matchedLabelPr = append(matchedLabelPr, gitapi.PrLabel{pr.Number, newLabels})
		}
	}

	return matchedLabelPr
}
