// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package server

import (
	"go.microcore.dev/framework/config/env"
)

const (
	// DefaultServiceName is the default name of the service.
	DefaultServiceName = "auth-service-server"
)

type (
	// Config defines the server configuration.
	Config struct {
		Name string
	}
)

// NewConfig creates and validates a server configuration.
func NewConfig() *Config {
	return &Config{
		Name: env.StrDefault("SERVICE_NAME", DefaultServiceName),
	}
}
