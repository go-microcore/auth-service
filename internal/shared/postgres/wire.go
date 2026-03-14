//go:build wireinject

// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package postgres

import (
	"fmt"

	"github.com/google/wire"
	"go.microcore.dev/auth-service/internal/shared/telemetry"
	"go.microcore.dev/framework/db/postgres"
	"go.microcore.dev/framework/db/postgres/client"
	"go.microcore.dev/framework/shutdown"
)

type (
	Options struct {
		Telemetry *telemetry.Telemetry
	}
)

func Init(opts *Options) (*Postgres, error) {
	wire.Build(
		NewConfig,
		newManager,
		newPostgres,
	)
	return nil, nil
}

func newManager(
	config *Config,
) (postgres.Manager, error) {
	manager, err := postgres.New(
		postgres.WithClientOptions(
			client.WithPostgresDSN(
				fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
					config.Host,
					config.Port,
					config.User,
					config.Password,
					config.DB,
					config.SSL,
				),
			),
		),
	)
	if err != nil {
		return nil, shutdown.NewExitReason(
			shutdown.ExitUnavailable,
			fmt.Errorf("create postgres manager: %w", err),
		)
	}
	return manager, nil
}

func newPostgres(
	config *Config,
	opts *Options,
	manager postgres.Manager,
) *Postgres {
	return &Postgres{
		Config:    config,
		Telemetry: opts.Telemetry,
		Manager:   manager,
	}
}
