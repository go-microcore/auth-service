//go:build wireinject

// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package telemetry

import (
	"context"
	"fmt"

	"github.com/google/wire"
	"go.microcore.dev/framework/shutdown"
	"go.microcore.dev/framework/telemetry"
)

func Init(ctx context.Context, opts *Options) (*Telemetry, error) {
	wire.Build(
		NewConfig,
		newManager,
		newTelemetry,
	)
	return nil, nil
}

func newManager(
	ctx context.Context,
	config *Config,
	opts *Options,
) (telemetry.Manager, error) {
	if config.Enabled {
		manager, err := telemetry.NewDefaultInsecureOtlpGrpc(
			ctx,
			config.Endpoint,
			opts.AppName,
		)
		if err != nil {
			return nil, shutdown.NewExitReason(
				shutdown.ExitUnavailable,
				fmt.Errorf("create telemetry manager: %w", err),
			)
		}
		return manager, nil
	}
	return nil, nil
}

func newTelemetry(
	options *Options,
	config *Config,
	manager telemetry.Manager,
) *Telemetry {
	return &Telemetry{
		Options: options,
		Config:  config,
		Manager: manager,
	}
}
