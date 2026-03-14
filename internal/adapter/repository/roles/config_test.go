// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package roles_test

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.microcore.dev/auth-service/internal/adapter/repository/roles"
)

func TestNewConfig(t *testing.T) {
	// Cache
	t.Setenv("CACHE_SUBSCRIBE_ROLE_UPDATES_BACKOFF_INITIAL", "2s")
	t.Setenv("CACHE_SUBSCRIBE_ROLE_UPDATES_BACKOFF_MAX", "20s")
	t.Setenv("CACHE_SUBSCRIBE_ROLE_UPDATES_BACKOFF_MULTIPLIER", "3")
	t.Setenv("CACHE_PERIODIC_ROLES_SYNC_INTERVAL", "2m")

	subscribeRoleUpdatesBackoffInitial, err := time.ParseDuration(
		os.Getenv("CACHE_SUBSCRIBE_ROLE_UPDATES_BACKOFF_INITIAL"),
	)
	require.NoError(t, err)

	subscribeRoleUpdatesBackoffMax, err := time.ParseDuration(
		os.Getenv("CACHE_SUBSCRIBE_ROLE_UPDATES_BACKOFF_MAX"),
	)
	require.NoError(t, err)

	subscribeRoleUpdatesBackoffMultiplier, err := strconv.ParseInt(
		os.Getenv("CACHE_SUBSCRIBE_ROLE_UPDATES_BACKOFF_MULTIPLIER"),
		10, 64,
	)
	require.NoError(t, err)

	periodicRolesSyncInterval, err := time.ParseDuration(
		os.Getenv("CACHE_PERIODIC_ROLES_SYNC_INTERVAL"),
	)
	require.NoError(t, err)

	cfg := roles.NewConfig()

	// Cache
	require.Equal(
		t,
		subscribeRoleUpdatesBackoffInitial,
		cfg.Cache.SubRoleUpBackoff.Initial,
	)
	require.Equal(
		t,
		subscribeRoleUpdatesBackoffMax,
		cfg.Cache.SubRoleUpBackoff.Max,
	)
	require.Equal(
		t,
		subscribeRoleUpdatesBackoffMultiplier,
		cfg.Cache.SubRoleUpBackoff.Multiplier,
	)
	require.Equal(
		t,
		periodicRolesSyncInterval,
		cfg.Cache.PeriodicRolesSyncInterval,
	)
}

//nolint:paralleltest // uses environment variables
func TestNewConfig_Default(t *testing.T) {
	cfg := roles.NewConfig()

	// Cache
	require.Equal(t,
		roles.DefaultCacheSubscribeRoleUpdatesBackoffInitial,
		cfg.Cache.SubRoleUpBackoff.Initial,
	)
	require.Equal(t,
		roles.DefaultCacheSubscribeRoleUpdatesBackoffMax,
		cfg.Cache.SubRoleUpBackoff.Max,
	)
	require.Equal(t,
		roles.DefaultCacheSubscribeRoleUpdatesBackoffMultiplier,
		cfg.Cache.SubRoleUpBackoff.Multiplier,
	)
	require.Equal(t,
		roles.DefaultCachePeriodicRolesSyncInterval,
		cfg.Cache.PeriodicRolesSyncInterval,
	)
}
