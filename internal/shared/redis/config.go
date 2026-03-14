// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package redis

import (
	"go.microcore.dev/framework/config/env"
)

const (
	// DefaultRedisAddr is the default Redis address.
	DefaultRedisAddr = "redis:6379"

	// DefaultRedisPassword is the default Redis password.
	DefaultRedisPassword = "password"

	// DefaultRedisDB is the default Redis database number.
	DefaultRedisDB = 0
)

type (
	// Config defines the Redis configuration.
	Config struct {
		Addr     string
		Password string
		DB       int
	}
)

// NewConfig creates and validates a Redis configuration.
func NewConfig() *Config {
	return &Config{
		Addr:     env.StrDefault("REDIS_ADDR", DefaultRedisAddr),
		Password: env.StrDefault("REDIS_PASSWORD", DefaultRedisPassword),
		DB:       env.IntDefault("REDIS_DB", DefaultRedisDB),
	}
}
