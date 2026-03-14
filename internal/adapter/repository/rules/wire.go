//go:build wireinject

// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package rules

import (
	"log/slog"

	"github.com/google/wire"
	"go.microcore.dev/auth-service/internal/port/adapter/repository/rules"
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

func Init(opts *Options) (rules.Adapter, error) {
	wire.Build(
		newLogger,
		newAdapter,
	)
	return nil, nil
}

func newLogger() *slog.Logger {
	return log.New("repository/rules")
}

func newAdapter(
	opts *Options,
	logger *slog.Logger,
) rules.Adapter {
	return NewAdapter(
		&AdapterConfig{
			Logger:   logger,
			Redis:    opts.Redis,
			Postgres: opts.Postgres,
		},
	)
}
