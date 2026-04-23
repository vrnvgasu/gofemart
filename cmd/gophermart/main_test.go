package main

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestParseConfig_Defaults(t *testing.T) {
	pflag.CommandLine = pflag.NewFlagSet("test", pflag.ContinueOnError)

	cfg := parseConfig()

	assert.NotNil(t, cfg)
	assert.Equal(t, ":8000", cfg.RunAddress)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, "secret", cfg.JWTSecret)
}
