// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package tokens

import (
	"time"

	"go.microcore.dev/framework/config/env"
)

const (
	// DefaultAuthMaxClockSkew is the default maximum allowed clock skew for auth requests.
	DefaultAuthMaxClockSkew = 5 * time.Second
)

type (
	// Config defines the tokens service configuration.
	Config struct {
		Auth *ConfigAuth
	}

	// ConfigAuth defines the auth configuration.
	ConfigAuth struct {
		MaxClockSkew time.Duration
	}
)

// NewConfig creates and validates a tokens service configuration.
func NewConfig() *Config {
	return &Config{
		Auth: &ConfigAuth{
			MaxClockSkew: env.DurDefault("AUTH_MAX_CLOCK_SKEW", DefaultAuthMaxClockSkew),
		},
	}
}
