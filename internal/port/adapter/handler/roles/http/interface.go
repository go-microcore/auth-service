// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package http

import (
	"context"

	"go.microcore.dev/framework/transport/http/server"
)

type (
	// Adapter defines the interface for the roles HTTP handler adapter.
	Adapter interface {
		CreateRole(
			ctx context.Context,
			reqCtx *server.RequestContext,
			req *CreateRoleRequest,
		)

		FilterRoles(
			ctx context.Context,
			reqCtx *server.RequestContext,
			req *FilterRolesRequest,
		)

		UpdateRole(
			ctx context.Context,
			reqCtx *server.RequestContext,
			req *UpdateRoleRequest,
		)

		DeleteRole(
			ctx context.Context,
			reqCtx *server.RequestContext,
		)
	}
)
