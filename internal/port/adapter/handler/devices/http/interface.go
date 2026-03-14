// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package http

import (
	"context"

	"go.microcore.dev/framework/transport/http/server"
)

type (
	// Adapter defines the interface for the devices HTTP handler adapter.
	Adapter interface {
		ListDevices(
			ctx context.Context,
			reqCtx *server.RequestContext,
		)
	}
)
