package config

import (
	"os"
	"testing"
)

func assertEqual(expectedValue interface{}, giventValue interface{}, t *testing.T) {
	if giventValue != expectedValue {
		t.Errorf("Expected value: %[1]s Found value: %[2]s", expectedValue, giventValue)
	}
}

func TestYamlConfigLoad(t *testing.T) {
	YamlPath = "./config_test.yaml"
	LoadYaml()
	t.Cleanup(func() {
		YamlConfig = YamlConfigV1{}
	})

	testData := []struct {
		config        string
		expectedValue string
	}{
		{YamlConfig.APIVersion, "v1"},
		{YamlConfig.Access.User, "tanmancan"},
		{YamlConfig.Access.Token, "testingtokenabcd"},
		{YamlConfig.Owner, "tanmancan"},
		{YamlConfig.Repo, "github-api-sandbox"},
	}

	for _, data := range testData {
		assertEqual(data.expectedValue, data.config, t)
	}

	t.Run("RuleConfigsExists", func(t *testing.T) {
		if len(YamlConfig.Rules) == 0 {
			t.Errorf("YamlConfig.Rules should not be empty")
		}

		rule := YamlConfig.Rules[0]
		label := "my-label-name"

		if rule.Label != label {
			t.Errorf("YamlConfig.Rules[0].Label should be %s", label)
		}

		ruleStringTestData := []struct {
			rule RuleTypeString
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
		{YamlConfig.APIVersion, ""},
		{YamlConfig.Access.User, ""},
		{YamlConfig.Access.Token, ""},
		{YamlConfig.Owner, ""},
		{YamlConfig.Repo, ""},
	}

	for _, data := range testData {
		assertEqual(data.expectedValue, data.config, t)
	}

	if len(YamlConfig.Rules) != 0 {
		t.Error("YamlConfig.Rules should be empty")
	}
}
func TestValidateVersion(t *testing.T) {
	APIVersion = "v2"
	validateVersion("v2")
}

func TestYamlGithubAccessUnmarshall(t *testing.T) {
	osToken := "TESTGITTOKEN"
	osUser := "TESTGITUSER"
	os.Setenv("GIT_TEST_TOKEN", osToken)
	os.Setenv("GIT_TEST_USER", osUser)

	YamlPath = "./config_test_env.yaml"
	LoadYaml()

	t.Cleanup(func() {
		YamlConfig = YamlConfigV1{}
		os.Unsetenv("GIT_TEST_TOKEN")
		os.Unsetenv("GIT_TEST_USER")
	})

	if YamlConfig.Access.Token != osToken {
		t.Errorf("YamlConfig.Access.Token should be %s", osToken)
	}

	if YamlConfig.Access.User != osUser {
		t.Errorf("YamlConfig.Access.User should be %s", osUser)
	}
}
