// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	tokenshp "go.microcore.dev/auth-service/internal/port/adapter/handler/tokens/http"
	rolessp "go.microcore.dev/auth-service/internal/port/service/roles"
	tokenssp "go.microcore.dev/auth-service/internal/port/service/tokens"
	"go.microcore.dev/framework/transport/http"
	"go.microcore.dev/framework/transport/http/server"
)

type (
	// AdapterConfig provides tokens HTTP adapter handler configuration.
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
func NewAdapter(config *AdapterConfig) tokenshp.Adapter {
	return &adapter{config}
}

// Auth issues JWT tokens.
//
// @Summary Issues JWT tokens.
// @Tags Tokens
// @Accept application/octet-stream
// @Produce application/octet-stream,json
// @Param request body tokenshp.AuthRequest true "BINARY DATA REQUIRED: Encrypt this JSON structure using AES-256-GCM with AUTH_KEY before sending"
// @Success 200 {object} tokenshp.AuthResponse "Returns encrypted binary. Decrypt using AES-256-GCM with AUTH_KEY to get this JSON structure."
// @Failure 400 {object} server.ErrResponse "Possible codes: INVALID_AUTH_REQUEST"
// @Failure 503 {object} server.ErrResponse "Possible codes: SERVICE_UNAVAILABLE"
// @Router /auth/tokens/ [post]
func (a *adapter) Auth(
	ctx context.Context,
	c *server.RequestContext,
) {
	// Decrypt request
	var req tokenshp.AuthRequest
	if err := a.TokensService.DecryptAuthRequest(
		ctx,
		c.Request.Body(),
		&req,
	); err != nil {
		c.WriteError(fmt.Errorf("decrypt request: %w", err))
		return
	}

	// Auth
	authData := tokenssp.AuthData(req)

	auth, err := a.TokensService.Auth(ctx, &authData)
	if err != nil {
		c.WriteError(fmt.Errorf("auth: %w", err))
		return
	}

	// Marshal raw response
	rawRes, err := json.Marshal(tokenshp.AuthResponse(*auth))
	if err != nil {
		c.WriteError(fmt.Errorf("json marshal response: %w", err))
		return
	}

	// Encrypt response
	encRes, err := a.TokensService.EncryptAuthResponse(
		ctx,
		rawRes,
	)
	if err != nil {
		c.WriteError(fmt.Errorf("encrypt response: %w", err))
		return
	}

	// Write res
	c.WriteWithStatusCode(http.StatusOK, encRes)
}

// Auth2fa issues JWT tokens after successful 2FA.
//
// @Summary Issues JWT tokens after successful 2FA.
// @Tags Tokens
// @Accept application/octet-stream
// @Produce application/octet-stream,json
// @Param request body tokenshp.Auth2FARequest true "BINARY DATA REQUIRED: Encrypt this JSON structure using AES-256-GCM with AUTH_KEY before sending"
// @Success 200 {object} tokenshp.Auth2FAResponse "Returns encrypted binary. Decrypt using AES-256-GCM with AUTH_KEY to get this JSON structure."
// @Failure 400 {object} server.ErrResponse "Possible codes: INVALID_AUTH_REQUEST"
// @Failure 503 {object} server.ErrResponse "Possible codes: SERVICE_UNAVAILABLE"
// @Router /auth/tokens/2fa [post]
func (a *adapter) Auth2fa(
	ctx context.Context,
	c *server.RequestContext,
) {
	// Decrypt request
	var req tokenshp.Auth2FARequest
	if err := a.TokensService.DecryptAuthRequest(
		ctx,
		c.Request.Body(),
		&req,
	); err != nil {
		c.WriteError(fmt.Errorf("decrypt request: %w", err))
		return
	}

	// Auth
	auth, err := a.TokensService.Auth2fa(
		ctx,
		tokenssp.Auth2faData(req),
	)
	if err != nil {
		c.WriteError(fmt.Errorf("auth: %w", err))
		return
	}

	// Marshal raw response
	rawRes, err := json.Marshal(tokenshp.Auth2FAResponse(*auth))
	if err != nil {
		c.WriteError(fmt.Errorf("json marshal response: %w", err))
		return
	}

	// Encrypt response
	encRes, err := a.TokensService.EncryptAuthResponse(
		ctx,
		rawRes,
	)
	if err != nil {
		c.WriteError(fmt.Errorf("encrypt response: %w", err))
		return
	}

	// Write res
	c.WriteWithStatusCode(http.StatusOK, encRes)
}

// TokenRenew renews JWT tokens based on a refresh token.
//
// @Summary Renews JWT tokens based on a refresh token.
// @Tags Tokens
// @Accept json
// @Produce json
// @Param request body tokenshp.TokenRenewRequest true "Request data"
// @Success 200 {object} tokenshp.TokenRenewResponse
// @Failure 400 {object} server.ErrResponse "Possible codes: INVALID_TOKEN, TOKEN_ALREADY_USED"
// @Failure 415 {object} server.ErrResponse "Possible codes: INVALID_JSON_BODY"
// @Failure 503 {object} server.ErrResponse "Possible codes: SERVICE_UNAVAILABLE"
// @Router /auth/tokens/renew [post]
func (a *adapter) TokenRenew(
	ctx context.Context,
	c *server.RequestContext,
	req *tokenshp.TokenRenewRequest,
) {
	// Renew token
	res, err := a.TokensService.TokenRenew(
		ctx,
		tokenssp.TokenRenewData(*req),
	)
	if err != nil {
		c.WriteError(fmt.Errorf("token renew: %w", err))
		return
	}

	// Write res
	c.WriteJsonWithStatusCode(
		http.StatusOK,
		tokenshp.TokenRenewResponse(*res),
	)
}

// TokenValidate validates an access token.
//
// @Summary Validates an access token.
// @Description Parsing and validation of the access token. Mainly used in the authorization mechanism across various microservices.
// @Tags Tokens
// @Security BearerAuth
// @Produce json
// @Success 200 {object} tokenshp.TokenValidateResponse
// @Failure 401 {object} server.ErrResponse "Possible codes: INVALID_TOKEN, 2FA_REQUIRED"
// @Failure 403 {object} server.ErrResponse "Possible codes: INSUFFICIENT_PERMISSIONS"
// @Failure 503 {object} server.ErrResponse "Possible codes: SERVICE_UNAVAILABLE"
// @Router /auth/tokens/validate [get]
func (a *adapter) TokenValidate(
	ctx context.Context,
	c *server.RequestContext,
) {
	// Get bearer token
	token, err := c.GetBearerToken()
	if err != nil {
		c.WriteError(fmt.Errorf("get bearer token: %w", err))
		return
	}

	// Validate access token
	res, err := a.TokensService.TokenValidate(
		ctx,
		tokenssp.TokenValidateData{
			AccessToken: token,
		},
	)
	if err != nil {
		c.WriteError(fmt.Errorf("token validate: %w", err))
		return
	}

	// Write res
	c.WriteJsonWithStatusCode(
		http.StatusOK,
		tokenshp.TokenValidateResponse(*res),
	)
}

// TokenAuthorizeHTTP checks HTTP authorization for given access token.
//
// @Summary Checks HTTP authorization for given access token.
// @Description Parsing, validation and authorize of the access token. Mainly used in the authorization mechanism across various microservices.
// @Tags Tokens
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body tokenshp.TokenAuthorizeHTTPRequest true "Request data"
// @Success 200 {object} tokenshp.TokenAuthorizeHTTPResponse
// @Failure 400 {object} server.ErrResponse "Possible codes: INVALID_PATH, INVALID_METHOD"
// @Failure 401 {object} server.ErrResponse "Possible codes: INVALID_TOKEN, 2FA_REQUIRED"
// @Failure 403 {object} server.ErrResponse "Possible codes: INSUFFICIENT_PERMISSIONS"
// @Failure 415 {object} server.ErrResponse "Possible codes: INVALID_JSON_BODY"
// @Failure 503 {object} server.ErrResponse "Possible codes: SERVICE_UNAVAILABLE"
// @Router /auth/tokens/authorize/http [post]
func (a *adapter) TokenAuthorizeHTTP(
	ctx context.Context,
	c *server.RequestContext,
	req *tokenshp.TokenAuthorizeHTTPRequest,
) {
	// Get bearer token
	token, err := c.GetBearerToken()
	if err != nil {
		c.WriteError(fmt.Errorf("get bearer token: %w", err))
		return
	}

	// Validate access token
	tokenData, err := a.TokensService.TokenValidate(
		ctx,
		tokenssp.TokenValidateData{
			AccessToken: token,
		},
	)
	if err != nil {
		c.WriteError(fmt.Errorf("token validate: %w", err))
		return
	}

	// Authorize http roles
	authData, err := a.RolesService.AuthorizeHTTPRoles(
		ctx,
		rolessp.AuthorizeHTTPRolesData{
			Roles:  tokenData.Roles,
			Path:   req.Path,
			Method: req.Method,
		},
	)
	if err != nil {
		c.WriteError(fmt.Errorf("authorize http roles: %w", err))
		return
	}

	// Write res
	c.WriteJsonWithStatusCode(
		http.StatusOK,
		tokenshp.TokenAuthorizeHTTPResponse{
			Token: tokenshp.TokenAuthorizeHTTPDataResponse(*tokenData),
			Auth:  tokenshp.TokenAuthorizeHTTPAuthResponse(*authData),
		},
	)
}

// Static access tokens

// CreateStaticAccessToken creates a static access token.
//
// @Summary Creates a static access token.
// @Tags Static tokens
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body tokenshp.CreateStaticAccessTokenRequest true "Request data"
// @Success 201 {object} tokenshp.CreateStaticAccessTokenResponse
// @Failure 400 {object} server.ErrResponse "Possible codes: INVALID_ID, INVALID_ROLES, INVALID_DESCRIPTION, STATIC_TOKEN_EXIST"
// @Failure 401 {object} server.ErrResponse "Possible codes: INVALID_TOKEN, 2FA_REQUIRED"
// @Failure 403 {object} server.ErrResponse "Possible codes: INSUFFICIENT_PERMISSIONS"
// @Failure 415 {object} server.ErrResponse "Possible codes: INVALID_JSON_BODY"
// @Failure 503 {object} server.ErrResponse "Possible codes: SERVICE_UNAVAILABLE"
// @Router /auth/tokens/static/ [post]
func (a *adapter) CreateStaticAccessToken(
	ctx context.Context,
	c *server.RequestContext,
	req *tokenshp.CreateStaticAccessTokenRequest,
) {
	// Create token
	res, err := a.TokensService.CreateStaticAccessToken(
		ctx,
		tokenssp.CreateStaticAccessTokenData(*req),
	)
	if err != nil {
		c.WriteError(fmt.Errorf("create token: %w", err))
		return
	}

	// Write res
	c.WriteJsonWithStatusCode(
		http.StatusCreated,
		tokenshp.CreateStaticAccessTokenResponse{
			Token: res,
		},
	)
}

// FilterStaticAccessTokens filters static access tokens.
//
// @Summary Filters static access tokens.
// @Tags Static tokens
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body tokenshp.FilterStaticAccessTokenRequest true "Request data"
// @Success 200 {array} tokenshp.FilterStaticAccessTokenResponse
// @Failure 401 {object} server.ErrResponse "Possible codes: INVALID_TOKEN, 2FA_REQUIRED"
// @Failure 403 {object} server.ErrResponse "Possible codes: INSUFFICIENT_PERMISSIONS"
// @Failure 415 {object} server.ErrResponse "Possible codes: INVALID_JSON_BODY"
// @Failure 503 {object} server.ErrResponse "Possible codes: SERVICE_UNAVAILABLE"
// @Router /auth/tokens/static/filter [post]
func (a *adapter) FilterStaticAccessTokens(
	ctx context.Context,
	c *server.RequestContext,
	req *tokenshp.FilterStaticAccessTokenRequest,
) {
	// Get tokens
	tokens, err := a.TokensService.FilterStaticAccessTokens(
		ctx,
		tokenssp.FilterStaticAccessTokenData(*req),
	)
	if err != nil {
		c.WriteError(fmt.Errorf("filter static access token: %w", err))
		return
	}

	// Make res
	res := make([]tokenshp.FilterStaticAccessTokenResponse, len(tokens))
	for i := range tokens {
		res[i] = tokenshp.FilterStaticAccessTokenResponse(tokens[i])
	}

	// Write res
	c.WriteJsonWithStatusCode(http.StatusOK, res)
}

// DeleteStaticAccessToken deletes a static access token by ID.
//
// @Summary Deletes a static access token by ID.
// @Tags Static tokens
// @Security BearerAuth
// @Produce json
// @Param id path string true "Token ID"
// @Success 204
// @Failure 400 {object} server.ErrResponse "Possible codes: INVALID_ID"
// @Failure 401 {object} server.ErrResponse "Possible codes: INVALID_TOKEN, 2FA_REQUIRED"
// @Failure 403 {object} server.ErrResponse "Possible codes: INSUFFICIENT_PERMISSIONS"
// @Failure 404 {object} server.ErrResponse "Possible codes: INVALID_ID, STATIC_TOKEN_NOT_FOUND"
// @Failure 503 {object} server.ErrResponse "Possible codes: SERVICE_UNAVAILABLE"
// @Router /auth/tokens/static/{id} [delete]
func (a *adapter) DeleteStaticAccessToken(
	ctx context.Context,
	c *server.RequestContext,
) {
	// Get token id
	id, err := c.UserValueStr("id")
	if err != nil {
		c.WriteError(tokenshp.ErrInvalidID)
		return
	}

	// Delete token
	if err := a.TokensService.DeleteStaticAccessToken(ctx, id); err != nil {
		c.WriteError(fmt.Errorf("delete static access token: %w", err))
		return
	}

	// Write res
	c.StatusCode(http.StatusNoContent)
}
