// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package redis

import (
	"fmt"

	"go.microcore.dev/auth-service/internal/shared/telemetry"
	"go.microcore.dev/framework/db/redis"
)

type (
	// Redis represents a Redis helper object.
	Redis struct {
		Config    *Config
		Telemetry *telemetry.Telemetry
		Manager   redis.Manager
	}
)

// Setup configures the helper object.
func (r *Redis) Setup() error {
	// Set telemetry manager
	if r.Telemetry.Config.Enabled {
		if err := r.Manager.SetTelemetryManager(r.Telemetry.Manager); err != nil {
			return fmt.Errorf("redis set telemetry: %w", err)
		}
	}

	return nil
}
