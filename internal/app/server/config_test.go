// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package server_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.microcore.dev/auth-service/internal/app/server"
)

func TestNewConfig(t *testing.T) {
	t.Setenv("SERVICE_NAME", "service")

	cfg := server.NewConfig()

	require.Equal(t, "service", cfg.Name)
}

//nolint:paralleltest // uses environment variables
func TestNewConfig_Default(t *testing.T) {
	cfg := server.NewConfig()

	require.Equal(t, server.DefaultServiceName, cfg.Name)
}
