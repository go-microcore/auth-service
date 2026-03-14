//go:build wireinject

// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package http

import (
	"log/slog"

	"github.com/google/wire"
	"go.microcore.dev/auth-service/internal/port/adapter/handler/roles/http"
	"go.microcore.dev/auth-service/internal/port/service/roles"
	"go.microcore.dev/framework/log"
)

type Options struct {
	RolesService roles.Service
}

func Init(opts *Options) (http.Adapter, error) {
	wire.Build(
		newLogger,
		newHandler,
	)
	return nil, nil
}

func newLogger() *slog.Logger {
	return log.New("handler/roles")
}

func newHandler(
	opts *Options,
	logger *slog.Logger,
) http.Adapter {
	return NewAdapter(
		&AdapterConfig{
			Logger:       logger,
			RolesService: opts.RolesService,
		},
	)
}
