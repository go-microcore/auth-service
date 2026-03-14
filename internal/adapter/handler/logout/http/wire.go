//go:build wireinject

// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package http

import (
	"log/slog"

	"github.com/google/wire"
	"go.microcore.dev/auth-service/internal/port/adapter/handler/logout/http"
	"go.microcore.dev/auth-service/internal/port/service/logout"
	"go.microcore.dev/framework/log"
)

type Options struct {
	LogoutService logout.Service
}

func Init(opts *Options) (http.Adapter, error) {
	wire.Build(
		newLogger,
		newHandler,
	)
	return nil, nil
}

func newLogger() *slog.Logger {
	return log.New("handler/logout")
}

func newHandler(
	opts *Options,
	logger *slog.Logger,
) http.Adapter {
	return NewAdapter(
		&AdapterConfig{
			Logger:        logger,
			LogoutService: opts.LogoutService,
		},
	)
}
