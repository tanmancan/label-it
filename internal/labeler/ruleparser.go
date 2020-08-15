package labeler

import (
	"regexp"
	"sort"
	"strconv"

	"github.com/tanmancan/label-it/v1/internal/common"
	"github.com/tanmancan/label-it/v1/internal/config"
	"github.com/tanmancan/label-it/v1/internal/gitapi"
)

// Rule label name and rules from YAML config
type Rule struct {
	Label       string
	HeadRules   config.RuleTypeString
	BaseRules   config.RuleTypeString
	TitleRules  config.RuleTypeString
	BodyRules   config.RuleTypeString
	UserRules   config.RuleTypeString
	NumberRules config.RuleTypeInt
	FileRules   config.RuleTypeString
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
	return RuleTypeStringValidator(r.UserRules, pr.User.Login)
}

// MatchNumberRules determines if pull request issue number matches provider number in rule
func (r Rule) MatchNumberRules(pr gitapi.PullRequest) bool {
	return RuleTypeIntValidator(r.NumberRules, pr.Number)
}

// MatchFileRules determines if changed files in pull request matches provided file rule
func (r Rule) MatchFileRules(pr gitapi.PullRequest) bool {
	rule := r.FileRules
	if rule == struct {
		Exact   string `yaml:"exact,omitempty"`
		NoExact string `yaml:"no-exact,omitempty"`
		Match   string `yaml:"match,omitempty"`
		NoMatch string `yaml:"no-match,omitempty"`
	}{} {
		return true
	}

	files := pr.Files
	valid := true
	invalid := false
	exact := valid
	noExact := valid
	match := valid
	noMatch := valid

	if rule.NoExact != "" {
		// No exact start as valid
		// But if a exact file is found, it will be marked invalid
		noExact = valid
		noExactIdx := sort.SearchStrings(files, rule.NoExact)

		// Exact match found
		if (noExactIdx != len(files)) && (files[noExactIdx] == rule.NoExact) {
			return false
		}
	}

	if rule.NoMatch != "" {
		// No match start as valid.
		// If a match is found using the no match check,
		// then no match will be marked as invalid
		noMatch = valid
		for _, file := range files {
			if matchString(rule.NoMatch, file) == true {
				noMatch = invalid
				break
			}
		}
		if noMatch == invalid {
			return false
		}
	}

	if rule.Exact != "" {
		// If exact rule is provided, we set default as invalid.
		// It will only pass if a match has been found
		exact = invalid
		exactIdx := sort.SearchStrings(files, rule.Exact)

		// No exact match found
		if exactIdx == len(files) {
			exact = invalid
		}

		// Exact match found
		if (exactIdx != len(files)) && (files[exactIdx] == rule.Exact) {
			exact = valid
		}
	}

	if rule.Match != "" {
		// If match check is provided,
		// we start it as invalid.
		// Match will only be flagged valid if a match is found.
		match = invalid
		for _, file := range files {
			if matchString(rule.Match, file) == true {
				match = valid
				break
			}
		}
	}

	return exact && noExact && match && noMatch
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
	case r.MatchFileRules(pr) != true:
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
	hasFileRule := false
	for _, rule := range config.YamlConfig.Rules {
		newRule := Rule{
			rule.Label,
			rule.Head,
			rule.Base,
			rule.Title,
			rule.Body,
			rule.User,
			rule.Number,
			rule.File,
		}
		labelRules = append(labelRules, newRule)
		if rule.File != struct {
			Exact   string `yaml:"exact,omitempty"`
			NoExact string `yaml:"no-exact,omitempty"`
			Match   string `yaml:"match,omitempty"`
			NoMatch string `yaml:"no-match,omitempty"`
		}{} {
			hasFileRule = true
		}
	}

	matchedLabelPr := []gitapi.PrLabel{}

	for _, pr := range prList {
		// Pre fetch files if file rule is present
		if hasFileRule == true {
			pr.Files = gitapi.GetAllFiles(pr.Number)
		}

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
