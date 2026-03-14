//go:build wireinject

// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package devices

import (
	"log/slog"

	"github.com/google/wire"
	authrp "go.microcore.dev/auth-service/internal/port/adapter/repository/auth"
	"go.microcore.dev/auth-service/internal/port/service/devices"
	"go.microcore.dev/framework/log"
)

type (
	Options struct {
		AuthRepository authrp.Adapter
	}
)

func Init(opts *Options) (devices.Service, error) {
	wire.Build(
		newLogger,
		newService,
	)
	return nil, nil
}

func newLogger() *slog.Logger {
	return log.New("service/devices")
}

func newService(
	opts *Options,
	logger *slog.Logger,
) devices.Service {
	return NewService(
		&ServiceConfig{
			Logger:         logger,
			AuthRepository: opts.AuthRepository,
		},
	)
}
