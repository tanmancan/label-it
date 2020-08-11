package config_test

import (
	"testing"

	"github.com/tanmancan/label-it/v1/internal/config"
)

func TestSetupArgs(t *testing.T) {
	err := config.SetupArgs()

	if err == nil {
		t.Errorf("SetupArgs should return an error if no config file was provided")
	}
}
