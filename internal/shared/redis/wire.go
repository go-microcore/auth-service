//go:build wireinject

// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package redis

import (
	"github.com/google/wire"
	"go.microcore.dev/auth-service/internal/shared/telemetry"
	"go.microcore.dev/framework/db/redis"
	"go.microcore.dev/framework/db/redis/client"
)

type (
	Options struct {
		Telemetry *telemetry.Telemetry
	}
)

func Init(opts *Options) (*Redis, error) {
	wire.Build(
		NewConfig,
		newManager,
		newRedis,
	)
	return nil, nil
}

func newManager(
	config *Config,
) redis.Manager {
	return redis.New(
		redis.WithClientOptions(
			client.WithAddr(config.Addr),
			client.WithPassword(config.Password),
			client.WithDB(config.DB),
		),
	)
}

func newRedis(
	config *Config,
	opts *Options,
	manager redis.Manager,
) *Redis {
	return &Redis{
		Config:    config,
		Telemetry: opts.Telemetry,
		Manager:   manager,
	}
}
