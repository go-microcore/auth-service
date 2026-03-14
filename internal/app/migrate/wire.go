//go:build wireinject

// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package migrate

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/wire"
	sharedpg "go.microcore.dev/auth-service/internal/shared/postgres"
	sharedtel "go.microcore.dev/auth-service/internal/shared/telemetry"
	"go.microcore.dev/framework/log"
)

func Init(ctx context.Context, opts *Options) (*Migrate, error) {
	wire.Build(
		NewConfig,
		newLogger,
		newTelemetry,
		newPostgres,
		newMigrate,
	)
	return nil, nil
}

func newLogger(
	cfg *Config,
) *slog.Logger {
	return log.New(cfg.Name)
}

func newTelemetry(
	ctx context.Context,
	cfg *Config,
) (*sharedtel.Telemetry, error) {
	t, err := sharedtel.Init(
		ctx,
		&sharedtel.Options{
			AppName: cfg.Name,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("init telemetry: %w", err)
	}
	return t, nil
}

func newPostgres(
	telemetry *sharedtel.Telemetry,
) (*sharedpg.Postgres, error) {
	p, err := sharedpg.Init(
		&sharedpg.Options{
			Telemetry: telemetry,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("init postgres: %w", err)
	}
	if err := p.Setup(); err != nil {
		return nil, fmt.Errorf("setup postgres: %w", err)
	}
	return p, nil
}

func newMigrate(
	config *Config,
	logger *slog.Logger,
	telemetry *sharedtel.Telemetry,
	postgres *sharedpg.Postgres,
) *Migrate {
	return &Migrate{
		Config:    config,
		Logger:    logger,
		Telemetry: telemetry,
		Postgres:  postgres,
	}
}
