// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package http

import (
	"go.microcore.dev/framework/transport/http/server"
)

type (
	// Adapter defines the interface for the auth HTTP handler adapter.
	Adapter interface {
		AuthMiddleware(
			handler server.RequestHandler,
		) server.RequestHandler
	}
)
