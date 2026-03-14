//go:build wireinject

// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package bootstrap

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/wire"
	sharedtel "go.microcore.dev/auth-service/internal/shared/telemetry"
	"go.microcore.dev/framework/log"
)

func Init(ctx context.Context, opts *Options) (*Bootstrap, error) {
	wire.Build(
		NewConfig,
		newLogger,
		newTelemetry,
		newBootstrap,
	)
	return nil, nil
}

func newLogger(
	cfg *Config,
) *slog.Logger {
	return log.New(cfg.Name)
}

func newTelemetry(
	ctx context.Context,
	cfg *Config,
) (*sharedtel.Telemetry, error) {
	t, err := sharedtel.Init(
		ctx,
		&sharedtel.Options{
			AppName: cfg.Name,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("init telemetry: %w", err)
	}
	return t, nil
}

func newBootstrap(
	options *Options,
	config *Config,
	logger *slog.Logger,
	telemetry *sharedtel.Telemetry,
) *Bootstrap {
	return &Bootstrap{
		Options:   options,
		Config:    config,
		Logger:    logger,
		Telemetry: telemetry,
	}
}
