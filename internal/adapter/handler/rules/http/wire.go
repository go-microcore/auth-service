//go:build wireinject

// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package http

import (
	"log/slog"

	"github.com/google/wire"
	"go.microcore.dev/auth-service/internal/port/adapter/handler/rules/http"
	"go.microcore.dev/auth-service/internal/port/service/rules"
	"go.microcore.dev/framework/log"
)

type Options struct {
	RulesService rules.Service
}

func Init(opts *Options) (http.Adapter, error) {
	wire.Build(
		newLogger,
		newHandler,
	)
	return nil, nil
}

func newLogger() *slog.Logger {
	return log.New("handler/rules")
}

func newHandler(
	opts *Options,
	logger *slog.Logger,
) http.Adapter {
	return NewAdapter(
		&AdapterConfig{
			Logger:       logger,
			RulesService: opts.RulesService,
		},
	)
}
