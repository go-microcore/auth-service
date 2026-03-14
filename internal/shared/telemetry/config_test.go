// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package telemetry_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.microcore.dev/auth-service/internal/shared/telemetry"
)

func TestNewConfig(t *testing.T) {
	t.Setenv("TELEMETRY_ENABLED", "true")
	t.Setenv("TELEMETRY_SERVER", "url")

	cfg := telemetry.NewConfig()

	require.True(t, cfg.Enabled)
	require.Equal(t, "url", cfg.Endpoint)
}

//nolint:paralleltest // uses environment variables
func TestNewConfig_Default(t *testing.T) {
	cfg := telemetry.NewConfig()

	require.Equal(t, telemetry.DefaultTelemetryEnabled, cfg.Enabled)
	require.Equal(t, telemetry.DefaultTelemetryServer, cfg.Endpoint)
}
