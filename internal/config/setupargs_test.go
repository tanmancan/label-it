package config_test

import (
	"testing"

	"github.com/tanmancan/label-it/v1/internal/config"
)

func TestSetupArgs(t *testing.T) {
	config.SetupArgs()
	tests := []struct {
		name string
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}
