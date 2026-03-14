//go:build wireinject

// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package roles

import (
	"log/slog"

	"github.com/google/wire"
	"go.microcore.dev/auth-service/internal/port/adapter/repository/roles"
	"go.microcore.dev/auth-service/internal/shared/postgres"
	"go.microcore.dev/auth-service/internal/shared/redis"
	"go.microcore.dev/framework/log"
)

type (
	Options struct {
		Postgres *postgres.Postgres
		Redis    *redis.Redis
	}
)

func Init(opts *Options) (roles.Adapter, error) {
	wire.Build(
		NewConfig,
		newLogger,
		newAdapter,
	)
	return nil, nil
}

func newLogger() *slog.Logger {
	return log.New("repository/roles")
}

func newAdapter(
	config *Config,
	opts *Options,
	logger *slog.Logger,
) roles.Adapter {
	return NewAdapter(
		&AdapterConfig{
			Config:   config,
			Logger:   logger,
			Redis:    opts.Redis,
			Postgres: opts.Postgres,
		},
	)
}
