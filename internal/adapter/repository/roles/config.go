// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package roles

import (
	"time"

	"go.microcore.dev/framework/config/env"
)

const (
	// DefaultCacheSubscribeRoleUpdatesBackoffInitial is the default initial backoff
	// duration for role updates subscription.
	DefaultCacheSubscribeRoleUpdatesBackoffInitial = 1 * time.Second

	// DefaultCacheSubscribeRoleUpdatesBackoffMax is the default maximum backoff
	// duration for role updates subscription.
	DefaultCacheSubscribeRoleUpdatesBackoffMax = 30 * time.Second

	// DefaultCacheSubscribeRoleUpdatesBackoffMultiplier is the default multiplier
	// for exponential backoff on role updates subscription.
	DefaultCacheSubscribeRoleUpdatesBackoffMultiplier = int64(2)

	// DefaultCachePeriodicRolesSyncInterval is the default interval for periodic
	// roles cache synchronization.
	DefaultCachePeriodicRolesSyncInterval = 5 * time.Minute
)

type (
	// Config defines roles configuration.
	Config struct {
		Cache *ConfigCache
	}

	// ConfigCache defines cache configuration.
	ConfigCache struct {
		SubRoleUpBackoff          *ConfigCacheSubRoleUpBackoff
		PeriodicRolesSyncInterval time.Duration
	}

	// ConfigCacheSubRoleUpBackoff defines subscribe role update backoff cache configuration.
	ConfigCacheSubRoleUpBackoff struct {
		Initial    time.Duration
		Max        time.Duration
		Multiplier int64
	}
)

// NewConfig creates and validates auth configuration.
func NewConfig() *Config {
	return &Config{
		Cache: &ConfigCache{
			SubRoleUpBackoff: &ConfigCacheSubRoleUpBackoff{
				Initial: env.DurDefault(
					"CACHE_SUBSCRIBE_ROLE_UPDATES_BACKOFF_INITIAL",
					DefaultCacheSubscribeRoleUpdatesBackoffInitial,
				),
				Max: env.DurDefault(
					"CACHE_SUBSCRIBE_ROLE_UPDATES_BACKOFF_MAX",
					DefaultCacheSubscribeRoleUpdatesBackoffMax,
				),
				Multiplier: env.Int64Default(
					"CACHE_SUBSCRIBE_ROLE_UPDATES_BACKOFF_MULTIPLIER",
					DefaultCacheSubscribeRoleUpdatesBackoffMultiplier,
				),
			},
			PeriodicRolesSyncInterval: env.DurDefault(
				"CACHE_PERIODIC_ROLES_SYNC_INTERVAL",
				DefaultCachePeriodicRolesSyncInterval,
			),
		},
	}
}
