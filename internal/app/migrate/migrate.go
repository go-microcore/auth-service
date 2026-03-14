// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package migrate

import (
	"context"
	"fmt"
	"log/slog"

	"go.microcore.dev/auth-service/internal/migrations"
	sharedpg "go.microcore.dev/auth-service/internal/shared/postgres"
	sharedtel "go.microcore.dev/auth-service/internal/shared/telemetry"
)

type (
	// Migrate represents a migrate application.
	Migrate struct {
		Config    *Config
		Logger    *slog.Logger
		Telemetry *sharedtel.Telemetry
		Postgres  *sharedpg.Postgres
	}

	// Options defines migrate primary options.
	Options struct{}
)

// Run executes the migrate application.
func (m *Migrate) Run(ctx context.Context) error {
	if err := m.Postgres.Manager.Migrate(migrations.Get(), nil); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	m.Logger.InfoContext(ctx, "migrations completed")

	return nil
}
