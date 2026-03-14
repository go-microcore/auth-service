// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package postgres_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	sharedpg "go.microcore.dev/auth-service/internal/shared/postgres"
)

func TestNewConfig(t *testing.T) {
	t.Setenv("POSTGRES_HOST", "host")
	t.Setenv("POSTGRES_PORT", "port")
	t.Setenv("POSTGRES_USER", "user")
	t.Setenv("POSTGRES_PASSWORD", "password")
	t.Setenv("POSTGRES_DB", "db")
	t.Setenv("POSTGRES_SSL", "ssl")

	cfg := sharedpg.NewConfig()

	require.Equal(t, "host", cfg.Host)
	require.Equal(t, "port", cfg.Port)
	require.Equal(t, "user", cfg.User)
	require.Equal(t, "password", cfg.Password)
	require.Equal(t, "db", cfg.DB)
	require.Equal(t, "ssl", cfg.SSL)
}

//nolint:paralleltest // uses environment variables
func TestNewConfig_Default(t *testing.T) {
	cfg := sharedpg.NewConfig()

	require.Equal(t, sharedpg.DefaultPostgresHost, cfg.Host)
	require.Equal(t, sharedpg.DefaultPostgresPort, cfg.Port)
	require.Equal(t, sharedpg.DefaultPostgresUser, cfg.User)
	require.Equal(t, sharedpg.DefaultPostgresPassword, cfg.Password)
	require.Equal(t, sharedpg.DefaultPostgresDB, cfg.DB)
	require.Equal(t, sharedpg.DefaultPostgresSSL, cfg.SSL)
}
