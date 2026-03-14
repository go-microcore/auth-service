// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package server

// @title		auth-service
// @version		1.0
// @BasePath	/

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

import (
	"context"
	"fmt"
	"log/slog"

	// Register Swagger docs init function.
	_ "go.microcore.dev/auth-service/docs"
	authhp "go.microcore.dev/auth-service/internal/port/adapter/handler/auth/http"
	deviceshp "go.microcore.dev/auth-service/internal/port/adapter/handler/devices/http"
	logouthp "go.microcore.dev/auth-service/internal/port/adapter/handler/logout/http"
	roleshp "go.microcore.dev/auth-service/internal/port/adapter/handler/roles/http"
	ruleshp "go.microcore.dev/auth-service/internal/port/adapter/handler/rules/http"
	tokenshp "go.microcore.dev/auth-service/internal/port/adapter/handler/tokens/http"
	rolesrp "go.microcore.dev/auth-service/internal/port/adapter/repository/roles"
	sharedhttp "go.microcore.dev/auth-service/internal/shared/http/server"
	sharedpg "go.microcore.dev/auth-service/internal/shared/postgres"
	sharedredis "go.microcore.dev/auth-service/internal/shared/redis"
	sharedtel "go.microcore.dev/auth-service/internal/shared/telemetry"
	"go.microcore.dev/framework/shutdown"
	"go.microcore.dev/framework/transport/http"
	"go.microcore.dev/framework/transport/http/server"
)

type (
	// Server represents a server application.
	Server struct {
		Config     *Config
		Logger     *slog.Logger
		Telemetry  *sharedtel.Telemetry
		Postgres   *sharedpg.Postgres
		Redis      *sharedredis.Redis
		HTTPServer *sharedhttp.Server
		Repository *Repository
		Handler    *Handler
	}

	// Repository wraps repository.
	Repository struct {
		Roles rolesrp.Adapter
	}

	// Handler wraps handlers.
	Handler struct {
		Auth    authhp.Adapter
		Devices deviceshp.Adapter
		Logout  logouthp.Adapter
		Roles   roleshp.Adapter
		Rules   ruleshp.Adapter
		Tokens  tokenshp.Adapter
	}

	// Options defines server primary options.
	Options struct{}
)

// Run executes the server application.
func (s *Server) Run(ctx context.Context) error {
	// Update roles cache
	if err := s.Repository.Roles.UpdateRolesCache(ctx); err != nil {
		return fmt.Errorf("update roles cache: %w", err)
	}

	// Subscribe role updates
	s.Repository.Roles.SubscribeRoleUpdates(ctx)

	// Start periodic roles cache sync
	s.Repository.Roles.PeriodicRolesCacheSync(ctx)

	// Setup routing
	s.SetupRouting()

	// Up http server
	go s.HTTPServer.Manager.Up()

	// Wait graceful shutdown
	code := shutdown.Wait()

	return shutdown.NewExitReason(code)
}

// SetupRouting registers all HTTP routes and groups for the server.
//
//nolint:funlen,dupl // need refactoring
func (s *Server) SetupRouting() {
	s.HTTPServer.Manager.
		// /auth
		AddRouteGroup(
			server.WithRouteGroupPath("/auth"),
			// /auth/tokens
			server.WithRouteGroup(
				server.WithRouteGroupPath("/tokens"),
				// /auth/tokens/static
				server.WithRouteGroup(
					server.WithRouteGroupPath("/static"),
					server.WithRouteGroupMiddlewares(
						s.Handler.Auth.AuthMiddleware,
					),
					// /auth/tokens/static/ [POST]
					server.WithRouteGroupRoute(
						server.WithRouteMethod(http.MethodPost),
						server.WithRouteBodyParserHandler(s.Handler.Tokens.CreateStaticAccessToken),
					),
					// /auth/tokens/static/filter [POST]
					server.WithRouteGroupRoute(
						server.WithRouteMethod(http.MethodPost),
						server.WithRoutePath("/filter"),
						server.WithRouteBodyParserHandler(s.Handler.Tokens.FilterStaticAccessTokens),
					),
					// /auth/tokens/static/{id} [DELETE]
					server.WithRouteGroupRoute(
						server.WithRouteMethod(http.MethodDelete),
						server.WithRoutePath("/{id}"),
						server.WithRouteHandler(s.Handler.Tokens.DeleteStaticAccessToken),
					),
				),
				// /auth/tokens/authorize
				server.WithRouteGroup(
					server.WithRouteGroupPath("/authorize"),
					server.WithRouteGroupMiddlewares(
						s.Handler.Auth.AuthMiddleware,
					),
					// /auth/tokens/authorize/http [POST]
					server.WithRouteGroupRoute(
						server.WithRouteMethod(http.MethodPost),
						server.WithRoutePath("/http"),
						server.WithRouteBodyParserHandler(s.Handler.Tokens.TokenAuthorizeHTTP),
					),
				),
				// /auth/tokens/ [POST]
				server.WithRouteGroupRoute(
					server.WithRouteMethod(http.MethodPost),
					server.WithRouteHandler(s.Handler.Tokens.Auth),
				),
				// /auth/tokens/2fa [POST]
				server.WithRouteGroupRoute(
					server.WithRouteMethod(http.MethodPost),
					server.WithRoutePath("/2fa"),
					server.WithRouteHandler(s.Handler.Tokens.Auth2fa),
				),
				// /auth/tokens/renew [POST]
				server.WithRouteGroupRoute(
					server.WithRouteMethod(http.MethodPost),
					server.WithRoutePath("/renew"),
					server.WithRouteBodyParserHandler(s.Handler.Tokens.TokenRenew),
				),
				// /auth/tokens/validate [GET]
				server.WithRouteGroupRoute(
					server.WithRouteMethod(http.MethodGet),
					server.WithRoutePath("/validate"),
					server.WithRouteHandler(s.Handler.Tokens.TokenValidate),
					server.WithRouteMiddlewares(
						s.Handler.Auth.AuthMiddleware,
					),
				),
			),
			// /auth/logout
			server.WithRouteGroup(
				server.WithRouteGroupPath("/logout"),
				server.WithRouteGroupMiddlewares(
					s.Handler.Auth.AuthMiddleware,
				),
				// /auth/logout/ [POST]
				server.WithRouteGroupRoute(
					server.WithRouteMethod(http.MethodPost),
					server.WithRouteHandler(s.Handler.Logout.Logout),
				),
				// /auth/logout/all [POST]
				server.WithRouteGroupRoute(
					server.WithRouteMethod(http.MethodPost),
					server.WithRoutePath("/all"),
					server.WithRouteHandler(s.Handler.Logout.LogoutAll),
				),
				// /auth/logout/device [POST]
				server.WithRouteGroupRoute(
					server.WithRouteMethod(http.MethodPost),
					server.WithRoutePath("/device"),
					server.WithRouteBodyParserHandler(s.Handler.Logout.LogoutDevice),
				),
			),
			// /auth/roles
			server.WithRouteGroup(
				server.WithRouteGroupPath("/roles"),
				server.WithRouteGroupMiddlewares(
					s.Handler.Auth.AuthMiddleware,
				),
				// /auth/roles/ [POST]
				server.WithRouteGroupRoute(
					server.WithRouteMethod(http.MethodPost),
					server.WithRouteBodyParserHandler(s.Handler.Roles.CreateRole),
				),
				// /auth/roles/filter [POST]
				server.WithRouteGroupRoute(
					server.WithRouteMethod(http.MethodPost),
					server.WithRoutePath("/filter"),
					server.WithRouteBodyParserHandler(s.Handler.Roles.FilterRoles),
				),
				// /auth/roles/{id} [PATCH]
				server.WithRouteGroupRoute(
					server.WithRouteMethod(http.MethodPatch),
					server.WithRoutePath("/{id}"),
					server.WithRouteBodyParserHandler(s.Handler.Roles.UpdateRole),
				),
				// /auth/roles/{id} [DELETE]
				server.WithRouteGroupRoute(
					server.WithRouteMethod(http.MethodDelete),
					server.WithRoutePath("/{id}"),
					server.WithRouteHandler(s.Handler.Roles.DeleteRole),
				),
			),
			// /auth/rules
			server.WithRouteGroup(
				server.WithRouteGroupPath("/rules"),
				server.WithRouteGroupMiddlewares(
					s.Handler.Auth.AuthMiddleware,
				),
				// /auth/rules/http
				server.WithRouteGroup(
					server.WithRouteGroupPath("/http"),
					// /auth/rules/http/ [POST]
					server.WithRouteGroupRoute(
						server.WithRouteMethod(http.MethodPost),
						server.WithRouteBodyParserHandler(s.Handler.Rules.CreateHTTPRule),
					),
					// /auth/rules/http/filter [POST]
					server.WithRouteGroupRoute(
						server.WithRouteMethod(http.MethodPost),
						server.WithRoutePath("/filter"),
						server.WithRouteBodyParserHandler(s.Handler.Rules.FilterHTTPRules),
					),
					// /auth/rules/http/{id} [PATCH]
					server.WithRouteGroupRoute(
						server.WithRouteMethod(http.MethodPatch),
						server.WithRoutePath("/{id}"),
						server.WithRouteBodyParserHandler(s.Handler.Rules.UpdateHTTPRule),
					),
					// /auth/rules/http/{id} [DELETE]
					server.WithRouteGroupRoute(
						server.WithRouteMethod(http.MethodDelete),
						server.WithRoutePath("/{id}"),
						server.WithRouteHandler(s.Handler.Rules.DeleteHTTPRule),
					),
				),
			),
			// /auth/devices [GET]
			server.WithRouteGroupRoute(
				server.WithRouteMethod(http.MethodGet),
				server.WithRoutePath("/devices"),
				server.WithRouteHandler(s.Handler.Devices.ListDevices),
				server.WithRouteMiddlewares(
					s.Handler.Auth.AuthMiddleware,
				),
			),
		)
}
