// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package postgres

import (
	"fmt"

	"go.microcore.dev/auth-service/internal/shared/telemetry"
	"go.microcore.dev/framework/db/postgres"
)

type (
	// Postgres represents a Postgres helper object.
	Postgres struct {
		Config    *Config
		Telemetry *telemetry.Telemetry
		Manager   postgres.Manager
	}
)

// Setup configures the helper object.
func (p *Postgres) Setup() error {
	// Set telemetry manager
	if p.Telemetry.Config.Enabled {
		if err := p.Manager.SetTelemetryManager(p.Telemetry.Manager); err != nil {
			return fmt.Errorf("postgres set telemetry: %w", err)
		}
	}

	return nil
}
