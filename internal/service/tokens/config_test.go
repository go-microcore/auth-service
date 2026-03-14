// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package tokens_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.microcore.dev/auth-service/internal/service/tokens"
)

func TestNewConfig(t *testing.T) {
	t.Setenv("AUTH_MAX_CLOCK_SKEW", "2s")

	authMaxClockSkew, err := time.ParseDuration(os.Getenv("AUTH_MAX_CLOCK_SKEW"))
	require.NoError(t, err)

	cfg := tokens.NewConfig()

	require.Equal(t, authMaxClockSkew, cfg.Auth.MaxClockSkew)
}

//nolint:paralleltest // uses environment variables
func TestNewConfig_Default(t *testing.T) {
	cfg := tokens.NewConfig()

	require.Equal(t, tokens.DefaultAuthMaxClockSkew, cfg.Auth.MaxClockSkew)
}
