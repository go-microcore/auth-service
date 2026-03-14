// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package http

import (
	"context"

	"go.microcore.dev/framework/transport/http/server"
)

type (
	// Adapter defines the interface for the tokens HTTP handler adapter.
	Adapter interface {
		Auth(
			ctx context.Context,
			reqCtx *server.RequestContext,
		)

		Auth2fa(
			ctx context.Context,
			reqCtx *server.RequestContext,
		)

		TokenRenew(
			ctx context.Context,
			reqCtx *server.RequestContext,
			req *TokenRenewRequest,
		)

		TokenValidate(
			ctx context.Context,
			reqCtx *server.RequestContext,
		)

		TokenAuthorizeHTTP(
			ctx context.Context,
			reqCtx *server.RequestContext,
			req *TokenAuthorizeHTTPRequest,
		)

		// Static access tokens

		CreateStaticAccessToken(
			ctx context.Context,
			reqCtx *server.RequestContext,
			req *CreateStaticAccessTokenRequest,
		)

		FilterStaticAccessTokens(
			ctx context.Context,
			reqCtx *server.RequestContext,
			req *FilterStaticAccessTokenRequest,
		)

		DeleteStaticAccessToken(
			ctx context.Context,
			reqCtx *server.RequestContext,
		)
	}
)
