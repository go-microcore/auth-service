//go:build wireinject

// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package tokens

import (
	"log/slog"

	"github.com/google/wire"
	authrp "go.microcore.dev/auth-service/internal/port/adapter/repository/auth"
	"go.microcore.dev/auth-service/internal/port/service/tokens"
	"go.microcore.dev/framework/log"
)

type (
	Options struct {
		AuthRepository authrp.Adapter
	}
)

func Init(opts *Options) (tokens.Service, error) {
	wire.Build(
		NewConfig,
		newLogger,
		newService,
	)
	return nil, nil
}

func newLogger() *slog.Logger {
	return log.New("service/tokens")
}

func newService(
	opts *Options,
	config *Config,
	logger *slog.Logger,
) tokens.Service {
	return NewService(
		&ServiceConfig{
			Config:         config,
			Logger:         logger,
			AuthRepository: opts.AuthRepository,
		},
	)
}
