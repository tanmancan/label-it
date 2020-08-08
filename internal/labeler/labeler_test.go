package labeler

import (
	"testing"

	"github.com/tanmancan/label-it/v1/internal/config"
)

func TestMatchString(t *testing.T) {
	pattern := "^(hello)"
	validString := "hello"
	invalidString := "world"

	match := matchString(pattern, validString)

	if match != true {
		t.Errorf("Pattern: %[1]s failed to match string '%[2]s'", pattern, validString)
	}

	noMatch := matchString(pattern, invalidString)

	if noMatch == true {
		t.Errorf("Pattern: %[1]s incorrectly matched string '%[2]s'", pattern, validString)
	}
}

func TestRuleTypeStringValidator(t *testing.T) {
	var valid bool
	var rule config.RuleTypeString

	t.Run("All Checks Provided", func(t *testing.T) {
		rule = config.RuleTypeString{
			Exact:   "octopus hello",
			NoExact: "world",
			Match:   "^(octopus)",
			NoMatch: "(lion)$",
		}

		compareString := "octopus hello"
		valid = RuleTypeStringValidator(rule, compareString)

		if valid == false {
			t.Error("RuleTypeString: Failed all checks")
		}
	})

	t.Run("Exact Only", func(t *testing.T) {
		rule = config.RuleTypeString{
			Exact: "AbcdEfg",
		}

		compareString := "AbcdEfg"
		valid = RuleTypeStringValidator(rule, compareString)

		if valid == false {
			t.Error("RuleTypeString: Failed exact check")
		}
	})

	t.Run("No Exact Only", func(t *testing.T) {
		rule = config.RuleTypeString{
			NoExact: "world",
		}

		compareString := "world"
		valid = RuleTypeStringValidator(rule, compareString)

		if valid != false {
			t.Error("RuleTypeString: Failed no exact check")
		}
	})

	t.Run("Match Only", func(t *testing.T) {
		rule = config.RuleTypeString{
			Match: "(more)",
		}

		compareString := "One or more words to match"
		valid = RuleTypeStringValidator(rule, compareString)

		if valid == false {
			t.Error("RuleTypeString: Failed match check")
		}
	})

	t.Run("No Match Only", func(t *testing.T) {
		rule = config.RuleTypeString{
			NoMatch: "(lion)$",
		}

		compareString := "King of the jungle, lion"
		valid = RuleTypeStringValidator(rule, compareString)

		if valid != false {
			t.Error("RuleTypeString: Failed no match check")
		}
	})
}
