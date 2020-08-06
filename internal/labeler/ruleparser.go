package labeler

import (
	"fmt"
	"label-it/internal/common"
	"label-it/internal/config"
	"label-it/internal/gitapi"
	"regexp"
	"sort"
	"strconv"
)

// Rule label name and rules from YAML config
type Rule struct {
	Label       string
	HeadRules   config.RuleTypeString
	BaseRules   config.RuleTypeString
	TitleRules  config.RuleTypeString
	BodyRules   config.RuleTypeString
	UserRule    config.RuleTypeString
	NumberRules config.RuleTypeInt
}

// LabelRules set of rules created from YAML config
type LabelRules []Rule

// Reusable function for pattern match
func matchString(pattern string, s string) bool {
	exp, experr := regexp.Compile(pattern)
	common.CheckErr(experr)

	return exp.MatchString(s)
}

// RuleTypeStringValidator validates a string value using rule group string
// Returns true if all rules validate, otherwise returns false
func RuleTypeStringValidator(r config.RuleTypeString, s string) bool {
	exact := r.Exact
	noExact := r.NoExact
	match := r.Match
	noMatch := r.NoMatch

	switch {
	case exact != "" && exact != s,
		noExact != "" && noExact == s,
		match != "" && matchString(match, s) != true,
		noMatch != "" && matchString(noMatch, s) == true:
		return false
	}

	return true
}

// RuleTypeIntValidator validates a int value using rule group integer
// Returns true if all rules validate, otherwise returns false
func RuleTypeIntValidator(r config.RuleTypeInt, i int) bool {
	exact := r.Exact
	noExact := r.NoExact
	match := r.Match
	noMatch := r.NoMatch

	// For regex pattern match, we coerce the int value into a string
	// @TODO: maybe just store the pr number as string and skip this step
	s := strconv.Itoa(i)

	switch {
	case exact != 0 && exact != i,
		noExact != 0 && noExact == i,
		match != "" && matchString(match, s) != true,
		noMatch != "" && matchString(noMatch, s) == true:
		return false
	}

	return true
}

// MatchHeadRules determines if provided pull request head branch matche the HeadRule
func (r Rule) MatchHeadRules(pr gitapi.PullRequest) bool {
	return RuleTypeStringValidator(r.HeadRules, pr.Head.Ref)
}

// MatchBaseRules determines if provided pull request base branch matche theBaseRule
func (r Rule) MatchBaseRules(pr gitapi.PullRequest) bool {
	return RuleTypeStringValidator(r.BaseRules, pr.Base.Ref)
}

// MatchTitleRules determines if provided pull request contains text in title rules
func (r Rule) MatchTitleRules(pr gitapi.PullRequest) bool {
	return RuleTypeStringValidator(r.TitleRules, pr.Title)
}

// MatchBodyRules determines if provided pull request contains text in title rules
func (r Rule) MatchBodyRules(pr gitapi.PullRequest) bool {
	return RuleTypeStringValidator(r.BodyRules, pr.Body)
}

// MatchUserRules checks if pull request creator username matches user rule
func (r Rule) MatchUserRules(pr gitapi.PullRequest) bool {
	return RuleTypeStringValidator(r.UserRule, pr.User.Login)
}

// MatchNumberRules determines if pull request issue number matches provider number in rule
func (r Rule) MatchNumberRules(pr gitapi.PullRequest) bool {
	return RuleTypeIntValidator(r.NumberRules, pr.Number)
}

func debugRules(r Rule, pr gitapi.PullRequest) {
	fmt.Println("----")
	fmt.Println(pr.Number, "Base:", pr.Base.Ref, r.MatchHeadRules(pr), r.BaseRules)
	fmt.Println(pr.Number, "Head:", pr.Head.Ref, r.MatchBaseRules(pr), r.HeadRules)
	fmt.Println(pr.Number, "Title", pr.Title, r.MatchTitleRules(pr), r.TitleRules)
	fmt.Println(pr.Number, "Body", pr.Body, r.MatchBodyRules(pr), r.BodyRules)
	fmt.Println(pr.Number, "User", pr.User.Login, r.MatchUserRules(pr), r.UserRule)
	fmt.Println(pr.Number, "Number", pr.Number, r.MatchNumberRules(pr), r.NumberRules)
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

	for _, prLabel := range pr.Labels {
		existingLabels = append(existingLabels, prLabel.Name)
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

	for _, rule := range config.YamlConfig.Rules {
		newRule := Rule{
			rule.Label,
			rule.Head,
			rule.Base,
			rule.Title,
			rule.Body,
			rule.User,
			rule.Number,
		}
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
			newPrLabel := gitapi.PrLabel{Issue: pr.Number, Labels: newLabels}
			matchedLabelPr = append(matchedLabelPr, newPrLabel)
		}
	}

	return matchedLabelPr
}
