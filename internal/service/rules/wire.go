//go:build wireinject

// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package rules

import (
	"log/slog"

	"github.com/google/wire"
	rulesrp "go.microcore.dev/auth-service/internal/port/adapter/repository/rules"
	"go.microcore.dev/auth-service/internal/port/service/rules"
	"go.microcore.dev/framework/log"
)

type (
	Options struct {
		RulesRepository rulesrp.Adapter
	}
)

func Init(opts *Options) (rules.Service, error) {
	wire.Build(
		newLogger,
		newService,
	)
	return nil, nil
}

func newLogger() *slog.Logger {
	return log.New("service/rules")
}

func newService(
	opts *Options,
	logger *slog.Logger,
) rules.Service {
	return NewService(
		&ServiceConfig{
			Logger:          logger,
			RulesRepository: opts.RulesRepository,
		},
	)
}
