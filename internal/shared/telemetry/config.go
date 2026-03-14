// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package telemetry

import (
	"go.microcore.dev/framework/config/env"
)

const (
	// DefaultTelemetryEnabled is the default value for enabling telemetry.
	DefaultTelemetryEnabled = false

	// DefaultTelemetryServer is the default telemetry server address.
	DefaultTelemetryServer = "alloy:4317"
)

type (
	// Config defines the telemetry configuration.
	Config struct {
		Enabled  bool
		Endpoint string
	}
)

// NewConfig creates and validates a telemetry configuration.
func NewConfig() *Config {
	return &Config{
		Enabled:  env.BoolDefault("TELEMETRY_ENABLED", DefaultTelemetryEnabled),
		Endpoint: env.StrDefault("TELEMETRY_SERVER", DefaultTelemetryServer),
	}
}
