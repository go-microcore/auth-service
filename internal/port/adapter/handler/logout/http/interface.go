// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package http

import (
	"context"

	"go.microcore.dev/framework/transport/http/server"
)

type (
	// Adapter defines the interface for the logout HTTP handler adapter.
	Adapter interface {
		Logout(
			ctx context.Context,
			reqCtx *server.RequestContext,
		)

		LogoutAll(
			ctx context.Context,
			reqCtx *server.RequestContext,
		)

		LogoutDevice(
			ctx context.Context,
			reqCtx *server.RequestContext,
			req *LogoutDeviceRequest,
		)
	}
)
