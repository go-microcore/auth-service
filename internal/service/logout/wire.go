//go:build wireinject

// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package logout

import (
	"log/slog"

	"github.com/google/wire"
	authrp "go.microcore.dev/auth-service/internal/port/adapter/repository/auth"
	"go.microcore.dev/auth-service/internal/port/service/logout"
	"go.microcore.dev/framework/log"
)

type (
	Options struct {
		AuthRepository authrp.Adapter
	}
)

func Init(opts *Options) (logout.Service, error) {
	wire.Build(
		newLogger,
		newService,
	)
	return nil, nil
}

func newLogger() *slog.Logger {
	return log.New("service/logout")
}

func newService(
	opts *Options,
	logger *slog.Logger,
) logout.Service {
	return NewService(
		&ServiceConfig{
			Logger:         logger,
			AuthRepository: opts.AuthRepository,
		},
	)
}
