//go:build wireinject

// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package server

import (
	"fmt"

	"github.com/google/wire"
	"go.microcore.dev/auth-service/internal/shared/telemetry"
	"go.microcore.dev/framework/shutdown"
	"go.microcore.dev/framework/transport/http/server"
	"go.microcore.dev/framework/transport/http/server/core"
	"go.microcore.dev/framework/transport/http/server/listener"
)

type (
	Options struct {
		Telemetry *telemetry.Telemetry
	}
)

func Init(opts *Options) (*Server, error) {
	wire.Build(
		NewConfig,
		newManager,
		newServer,
	)
	return nil, nil
}

func newManager(
	config *Config,
) (server.Manager, error) {
	manager, err := server.New(
		server.WithListenerOptions(
			listener.WithHostname(config.Host),
			listener.WithPort(config.Port),
		),
		server.WithCoreOptions(
			core.WithName(config.Name),
			core.WithConcurrency(config.Concurrency),
			core.WithReadBufferSize(config.ReadBufferSize),
			core.WithWriteBufferSize(config.WriteBufferSize),
			core.WithReadTimeout(config.ReadTimeout),
			core.WithWriteTimeout(config.WriteTimeout),
			core.WithIdleTimeout(config.IdleTimeout),
			core.WithMaxConnsPerIP(config.MaxConnsPerIP),
			core.WithMaxRequestsPerConn(config.MaxRequestsPerConn),
			core.WithMaxRequestBodySize(config.MaxRequestBodySize),
			core.WithDisableKeepalive(config.DisableKeepalive),
			core.WithTCPKeepalive(config.TCPKeepalive),
			core.WithLogAllErrors(config.LogAllErrors),
		),
	)
	if err != nil {
		return nil, shutdown.NewExitReason(
			shutdown.ExitOSError,
			fmt.Errorf("create http server manager: %w", err),
		)
	}
	return manager, nil
}

func newServer(
	config *Config,
	opts *Options,
	manager server.Manager,
) *Server {
	return &Server{
		Config:    config,
		Telemetry: opts.Telemetry,
		Manager:   manager,
	}
}
