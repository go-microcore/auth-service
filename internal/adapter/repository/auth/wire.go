//go:build wireinject

// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package auth

import (
	"fmt"
	"log/slog"

	"github.com/google/wire"
	"go.microcore.dev/auth-service/internal/port/adapter/repository/auth"
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

func Init(opts *Options) (auth.Adapter, error) {
	wire.Build(
		NewConfig,
		newLogger,
		newLocalTokenCache,
		newAdapter,
	)
	return nil, nil
}

func newLogger() *slog.Logger {
	return log.New("repository/auth")
}

func newLocalTokenCache(
	config *Config,
) (*LocalTokenCache, error) {
	cache, err := NewLocalTokenCache(config.Cache.LocalTokenSize)
	if err != nil {
		return nil, fmt.Errorf("local token cache: %w", err)
	}
	return cache, nil
}

func newAdapter(
	opts *Options,
	config *Config,
	logger *slog.Logger,
	localTokenCache *LocalTokenCache,
) auth.Adapter {
	return NewAdapter(
		&AdapterConfig{
			Config:          config,
			Logger:          logger,
			Redis:           opts.Redis,
			Postgres:        opts.Postgres,
			LocalTokenCache: localTokenCache,
		},
	)
}
