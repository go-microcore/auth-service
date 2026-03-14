// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package http

import (
	"context"

	"go.microcore.dev/framework/transport/http/server"
)

type (
	// Adapter defines the interface for the rules HTTP handler adapter.
	Adapter interface {
		CreateHTTPRule(
			ctx context.Context,
			reqCtx *server.RequestContext,
			req *CreateHTTPRuleRequest,
		)

		FilterHTTPRules(
			ctx context.Context,
			reqCtx *server.RequestContext,
			req *FilterHTTPRulesRequest,
		)

		UpdateHTTPRule(
			ctx context.Context,
			reqCtx *server.RequestContext,
			req *UpdateHTTPRuleRequest,
		)

		DeleteHTTPRule(
			ctx context.Context,
			reqCtx *server.RequestContext,
		)
	}
)
