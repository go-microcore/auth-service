//go:build wireinject

// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package roles

import (
	"log/slog"

	"github.com/google/wire"
	rolesrp "go.microcore.dev/auth-service/internal/port/adapter/repository/roles"
	"go.microcore.dev/auth-service/internal/port/service/roles"
	"go.microcore.dev/framework/log"
)

type (
	Options struct {
		RolesRepository rolesrp.Adapter
	}
)

func Init(opts *Options) (roles.Service, error) {
	wire.Build(
		newLogger,
		newService,
	)
	return nil, nil
}

func newLogger() *slog.Logger {
	return log.New("service/roles")
}

func newService(
	opts *Options,
	logger *slog.Logger,
) roles.Service {
	return NewService(
		&ServiceConfig{
			Logger:          logger,
			RolesRepository: opts.RolesRepository,
		},
	)
}
