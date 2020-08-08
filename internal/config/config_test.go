package config

import (
	"os"
	"testing"
)

func TestYamlConfigLoad(t *testing.T) {
	YamlPath = "./config_test.yaml"
	LoadYaml()
	t.Cleanup(func() {
		YamlConfig = YamlConfigV1{}
	})

	t.Run("ApiConfigExists", func(t *testing.T) {
		ver := "v1"
		if YamlConfig.APIVersion != ver {
			t.Errorf("YamlConfig.APIVersion should be %s", ver)
		}
		user := "tanmancan"
		if YamlConfig.Access.User != user {
			t.Errorf("YamlConfig.Access.User should be %s", user)
		}
		token := "testingtokenabcd"
		if YamlConfig.Access.Token != token {
			t.Errorf("YamlConfig.Access.User should be %s", token)
		}
		owner := "tanmancan"
		if YamlConfig.Owner != owner {
			t.Errorf("YamlConfig.Owner should be %s", owner)
		}
		repo := "github-api-sandbox"
		if YamlConfig.Repo != "github-api-sandbox" {
			t.Errorf("YamlConfig.Repo should be %s", repo)
		}
	})

	t.Run("RuleConfigsExists", func(t *testing.T) {
		if len(YamlConfig.Rules) == 0 {
			t.Errorf("YamlConfig.Rules should not be empty")
		}

		rule := YamlConfig.Rules[0]
		label := "my-label-name"

		if rule.Label != label {
			t.Errorf("YamlConfig.Rules[0].Label should be %s", label)
		}

		headRule := YamlConfig.Rules[0].Head
		exact := "master"
		noExact := "staging"
		match := "^(mas)"
		noMatch := "^(stag)"

		if headRule.Exact != exact {
			t.Errorf("YamlConfig.Rules[0].Head.Exact must be %s", exact)
		}
		if headRule.NoExact != noExact {
			t.Errorf("YamlConfig.Rules[0].Head.NoExact must be %s", noExact)
		}
		if headRule.Match != match {
			t.Errorf("YamlConfig.Rules[0].Head.Match must be %s", match)
		}
		if headRule.NoMatch != noMatch {
			t.Errorf("YamlConfig.Rules[0].Head.NoMatch must be %s", noMatch)
		}
	})
}

func TestYamlConfigInitial(t *testing.T) {
	if YamlConfig.APIVersion != "" {
		t.Error("YamlConfig.APIVersion should be empty")
	}
	if YamlConfig.Access.User != "" {
		t.Error("YamlConfig.Access.User should be empty")
	}
	if YamlConfig.Access.Token != "" {
		t.Error("YamlConfig.Access.User should be empty")
	}
	if YamlConfig.Owner != "" {
		t.Error("YamlConfig.Owner should be empty")
	}
	if YamlConfig.Repo != "" {
		t.Error("YamlConfig.Repo should be empty")
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
