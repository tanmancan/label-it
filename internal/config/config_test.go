package config_test

import (
	"os"
	"testing"

	"github.com/tanmancan/label-it/v1/internal/config"
)

func assertEqual(expectedValue interface{}, givenValue interface{}, t *testing.T) {
	if givenValue != expectedValue {
		t.Errorf("Expected value: %[1]s Found value: %[2]s", expectedValue, givenValue)
	}
}

func TestYamlConfigLoad(t *testing.T) {
	config.YamlPath = "./config_test.yaml"
	config.LoadYaml()
	t.Cleanup(func() {
		config.YamlConfig = config.YamlConfigV1{}
	})

	testData := []struct {
		config        string
		expectedValue string
	}{
		{config.YamlConfig.APIVersion, "v1"},
		{config.YamlConfig.Access.User, "tanmancan"},
		{config.YamlConfig.Access.Token, "testingTokenAbcd"},
		{config.YamlConfig.Owner, "tanmancan"},
		{config.YamlConfig.Repo, "github-api-sandbox"},
	}

	for _, data := range testData {
		assertEqual(data.expectedValue, data.config, t)
	}

	t.Run("RuleConfigsExists", func(t *testing.T) {
		if len(config.YamlConfig.Rules) == 0 {
			t.Errorf("config.YamlConfig.Rules should not be empty")
		}

		rule := config.YamlConfig.Rules[0]
		label := "my-label-name"

		if rule.Label != label {
			t.Errorf("config.YamlConfig.Rules[0].Label should be %s", label)
		}

		ruleStringTestData := []struct {
			rule config.RuleTypeString
		}{
			{rule.Head},
			{rule.Base},
			{rule.Title},
			{rule.Body},
			{rule.User},
		}

		exact := "master"
		noExact := "staging"
		match := "^(mas)"
		noMatch := "^(stag)"

		for _, stringRule := range ruleStringTestData {
			t.Run("ExactCheck", func(t *testing.T) {
				assertEqual(exact, stringRule.rule.Exact, t)
			})
			t.Run("NoExactCheck", func(t *testing.T) {
				assertEqual(noExact, stringRule.rule.NoExact, t)
			})
			t.Run("MatchCheck", func(t *testing.T) {
				assertEqual(match, stringRule.rule.Match, t)
			})
			t.Run("NoMatchCheck", func(t *testing.T) {
				assertEqual(noMatch, stringRule.rule.NoMatch, t)
			})
		}

	})
}

func TestYamlConfigInitial(t *testing.T) {
	testData := []struct {
		config        string
		expectedValue string
	}{
		{config.YamlConfig.APIVersion, ""},
		{config.YamlConfig.Access.User, ""},
		{config.YamlConfig.Access.Token, ""},
		{config.YamlConfig.Owner, ""},
		{config.YamlConfig.Repo, ""},
	}

	for _, data := range testData {
		assertEqual(data.expectedValue, data.config, t)
	}

	if len(config.YamlConfig.Rules) != 0 {
		t.Error("config.YamlConfig.Rules should be empty")
	}
}

func TestYamlGithubAccessUnmarshal(t *testing.T) {
	osToken := "TESTGITTOKEN"
	osUser := "TESTGITUSER"
	os.Setenv("GIT_TEST_TOKEN", osToken)
	os.Setenv("GIT_TEST_USER", osUser)

	config.YamlPath = "./config_test_env.yaml"
	config.LoadYaml()

	t.Cleanup(func() {
		config.YamlConfig = config.YamlConfigV1{}
		os.Unsetenv("GIT_TEST_TOKEN")
		os.Unsetenv("GIT_TEST_USER")
	})

	if config.YamlConfig.Access.Token != osToken {
		t.Errorf("config.YamlConfig.Access.Token should be %s", osToken)
	}

	if config.YamlConfig.Access.User != osUser {
		t.Errorf("config.YamlConfig.Access.User should be %s", osUser)
	}
}
