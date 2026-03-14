// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package bootstrap

import (
	"go.microcore.dev/framework/config/env"
)

const (
	// DefaultServiceName is the default name of the service.
	DefaultServiceName = "auth-service-bootstrap"
)

type (
	// Config defines the bootstrap configuration.
	Config struct {
		Name string
	}
)

// NewConfig creates and validates a bootstrap configuration.
func NewConfig() *Config {
	return &Config{
		Name: env.StrDefault("SERVICE_NAME", DefaultServiceName),
	}
}
