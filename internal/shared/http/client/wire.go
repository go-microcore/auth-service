//go:build wireinject

// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package client

import (
	"github.com/google/wire"
	"go.microcore.dev/auth-service/internal/shared/telemetry"
	"go.microcore.dev/framework/transport/http/client"
)

type (
	Options struct {
		Telemetry *telemetry.Telemetry
	}
)

func Init(opts *Options) (*Client, error) {
	wire.Build(
		newManager,
		newClient,
	)
	return nil, nil
}

func newManager() client.Manager {
	return client.New()
}

func newClient(
	opts *Options,
	manager client.Manager,
) *Client {
	return &Client{
		Telemetry: opts.Telemetry,
		Manager:   manager,
	}
}
