// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package http

import (
	"fmt"
	"log/slog"

	authhp "go.microcore.dev/auth-service/internal/port/adapter/handler/auth/http"
	rolessp "go.microcore.dev/auth-service/internal/port/service/roles"
	tokenssp "go.microcore.dev/auth-service/internal/port/service/tokens"
	"go.microcore.dev/framework/transport/http/server"
)

type (
	// AdapterConfig provides auth HTTP adapter handler configuration.
	AdapterConfig struct {
		Logger        *slog.Logger
		TokensService tokenssp.Service
		RolesService  rolessp.Service
	}

	adapter struct {
		*AdapterConfig
	}
)

// NewAdapter creates a new instance of the adapter.
func NewAdapter(config *AdapterConfig) authhp.Adapter {
	return &adapter{config}
}

// AuthMiddleware handles HTTP request role authorization.
func (a *adapter) AuthMiddleware(
	handler server.RequestHandler,
) server.RequestHandler {
	return func(c *server.RequestContext) {
		tokenData, err := a.validateAccessToken(c)
		if err != nil {
			c.WriteError(err)
			return
		}

		authData, err := a.authorizeHTTP(c, tokenData.Roles)
		if err != nil {
			c.WriteError(err)
			return
		}

		if authData.Mfa && tokenData.Mfa {
			c.WriteError(authhp.ErrMFARequired)
			return
		}

		setAuthContext(c, tokenData, authData)
		handler(c)
	}
}

func (a *adapter) validateAccessToken(
	c *server.RequestContext,
) (*tokenssp.TokenValidateResult, error) {
	token, err := c.GetBearerToken()
	if err != nil {
		return nil, authhp.ErrInvalidToken
	}

	res, err := a.TokensService.TokenValidate(
		c.GetContext(),
		tokenssp.TokenValidateData{
			AccessToken: token,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("validate token: %w", err)
	}

	return res, nil
}

func (a *adapter) authorizeHTTP(
	c *server.RequestContext,
	roles []string,
) (*rolessp.AuthorizeHTTPRolesResult, error) {
	res, err := a.RolesService.AuthorizeHTTPRoles(
		c.GetContext(),
		rolessp.AuthorizeHTTPRolesData{
			Roles:  roles,
			Path:   string(c.Path()),
			Method: string(c.Method()),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("authorize http roles: %w", err)
	}

	return res, nil
}

func setAuthContext(
	c *server.RequestContext,
	tokenData *tokenssp.TokenValidateResult,
	authData *rolessp.AuthorizeHTTPRolesResult,
) {
	c.SetUserValue("device", tokenData.Device)
	c.SetUserValue("user", tokenData.User)
	c.SetUserValue("roles", tokenData.Roles)
	c.SetUserValue("mfa_value", tokenData.Mfa)
	c.SetUserValue("mfa_validation", authData.Mfa)
}
