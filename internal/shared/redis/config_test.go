// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package redis_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	sharedredis "go.microcore.dev/auth-service/internal/shared/redis"
)

func TestNewConfig(t *testing.T) {
	t.Setenv("REDIS_ADDR", "addr")
	t.Setenv("REDIS_PASSWORD", "password")
	t.Setenv("REDIS_DB", "0")

	dbExpected := 0

	cfg := sharedredis.NewConfig()

	require.Equal(t, "addr", cfg.Addr)
	require.Equal(t, "password", cfg.Password)
	require.Equal(t, dbExpected, cfg.DB)
}

//nolint:paralleltest // uses environment variables
func TestNewConfig_Default(t *testing.T) {
	cfg := sharedredis.NewConfig()

	require.Equal(t, sharedredis.DefaultRedisAddr, cfg.Addr)
	require.Equal(t, sharedredis.DefaultRedisPassword, cfg.Password)
	require.Equal(t, sharedredis.DefaultRedisDB, cfg.DB)
}
