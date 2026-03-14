// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package tokens

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	authrp "go.microcore.dev/auth-service/internal/port/adapter/repository/auth"
	tokenssp "go.microcore.dev/auth-service/internal/port/service/tokens"
	"go.microcore.dev/framework/transport"
)

type (
	// ServiceConfig provides tokens service configuration.
	ServiceConfig struct {
		Config         *Config
		Logger         *slog.Logger
		AuthRepository authrp.Adapter
	}

	service struct {
		*ServiceConfig
	}
)

// NewService creates a new instance of the service.
func NewService(config *ServiceConfig) tokenssp.Service {
	return &service{config}
}

// DecryptAuthRequest decrypt and parse auth request.
func (s *service) DecryptAuthRequest(
	ctx context.Context,
	data []byte,
	target any,
) error {
	// Decrypt raw request
	rawReq, err := s.AuthRepository.DecryptAuthRequest(ctx, data)
	if err != nil {
		s.Logger.ErrorContext(
			ctx,
			"decrypt auth request",
			slog.Any("error", err),
			slog.Any("data", data),
		)

		return tokenssp.ErrInvalidAuthRequest
	}

	// Unmarshal raw request
	if err := json.Unmarshal(rawReq, target); err != nil {
		s.Logger.ErrorContext(
			ctx,
			"json unmarshal auth request",
			slog.Any("error", err),
			slog.Any("data", data),
		)

		return tokenssp.ErrInvalidAuthRequest
	}

	return nil
}

// EncryptAuthResponse encrypt auth response.
func (s *service) EncryptAuthResponse(
	ctx context.Context,
	data []byte,
) ([]byte, error) {
	// Encrypt response
	res, err := s.AuthRepository.EncryptAuthResponse(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("encrypt: %w", err)
	}

	return res, nil
}

// Auth issues JWT tokens.
func (s *service) Auth(
	ctx context.Context,
	data *tokenssp.AuthData,
) (*tokenssp.AuthResult, error) {
	// Validate TTL
	if err := s.validateTTL(data.TTL); err != nil {
		return nil, err
	}

	// Gen new tokens
	tokens, err := s.AuthRepository.NewTokens(
		ctx,
		authrp.NewTokenData{
			User:   data.User,
			Roles:  data.Roles,
			Mfa:    data.Mfa,
			Device: data.Device,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("new tokens: %w", err)
	}

	// Parse refresh token
	refreshToken, err := s.AuthRepository.ParseRefreshToken(ctx, tokens.Refresh)
	if err != nil {
		return nil, fmt.Errorf("parse refresh token: %w", err)
	}

	// Upsert session
	newDevice, err := s.upsertSession(ctx, data, refreshToken.ID)
	if err != nil {
		return nil, err
	}

	// Return tokens
	return &tokenssp.AuthResult{
		Access:    tokens.Access,
		Refresh:   tokens.Refresh,
		Mfa:       data.Mfa,
		NewDevice: newDevice,
	}, nil
}

// Auth2fa issues JWT tokens after successful 2FA.
func (s *service) Auth2fa(
	ctx context.Context,
	data tokenssp.Auth2faData,
) (*tokenssp.Auth2faResult, error) {
	// Validate TTL
	if err := s.validateTTL(data.TTL); err != nil {
		return nil, err
	}

	// Gen new tokens
	tokens, err := s.AuthRepository.NewTokens(
		ctx,
		authrp.NewTokenData{
			User:   data.User,
			Roles:  data.Roles,
			Mfa:    false,
			Device: data.Device,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("new tokens: %w", err)
	}

	// Parse refresh token
	refreshToken, err := s.AuthRepository.ParseRefreshToken(ctx, tokens.Refresh)
	if err != nil {
		return nil, fmt.Errorf("parse refresh token: %w", err)
	}

	// Update session
	if err := s.AuthRepository.UpdateSession(
		ctx,
		data.User,
		data.Device,
		refreshToken.ID,
	); err != nil {
		return nil, fmt.Errorf("update session: %w", err)
	}

	// Return tokens
	result := tokenssp.Auth2faResult(*tokens)

	return &result, nil
}

// TokenRenew renews JWT tokens based on a refresh token.
func (s *service) TokenRenew(
	ctx context.Context,
	data tokenssp.TokenRenewData,
) (*tokenssp.TokenRenewResult, error) {
	// Parse refresh token
	token, err := s.AuthRepository.ParseRefreshToken(ctx, data.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("parse refresh token: %w", err)
	}

	// Get session
	session, err := s.AuthRepository.GetSession(ctx, token.User, token.Device)
	if err != nil && !errors.Is(err, authrp.ErrSessionNotFound) {
		return nil, fmt.Errorf("get session: %w", err)
	}

	// Check session exist
	if session == nil {
		return nil, tokenssp.ErrRefreshTokenAlreadyUsed
	}

	// Check refresh token jti
	if session.Jti != token.ID {
		return nil, tokenssp.ErrRefreshTokenAlreadyUsed
	}

	// Create tokens
	tokens, err := s.AuthRepository.NewTokens(
		ctx,
		authrp.NewTokenData{
			User:   token.User,
			Roles:  token.Roles,
			Mfa:    token.Mfa,
			Device: token.Device,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("new tokens: %w", err)
	}

	// Parse refresh token
	refreshToken, err := s.AuthRepository.ParseRefreshToken(ctx, tokens.Refresh)
	if err != nil {
		return nil, fmt.Errorf("parse new refresh token: %w", err)
	}

	// Update session
	if err := s.AuthRepository.UpdateSession(
		ctx,
		refreshToken.User,
		token.Device,
		refreshToken.ID,
	); err != nil {
		return nil, fmt.Errorf("update session: %w", err)
	}

	// Return tokens and MFA flag
	return &tokenssp.TokenRenewResult{
		Access:  tokens.Access,
		Refresh: tokens.Refresh,
		Mfa:     token.Mfa,
	}, nil
}

// TokenValidate validates an access token.
func (s *service) TokenValidate(
	ctx context.Context,
	data tokenssp.TokenValidateData,
) (*tokenssp.TokenValidateResult, error) {
	// Parse and validate access token
	t, err := s.AuthRepository.ParseAccessToken(ctx, data.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	// Make res
	res := tokenssp.TokenValidateResult(*t)

	return &res, nil
}

// Static access tokens

// CreateStaticAccessToken creates a static access token.
func (s *service) CreateStaticAccessToken(
	ctx context.Context,
	data tokenssp.CreateStaticAccessTokenData,
) (string, error) {
	// Check token exist
	if tokens, err := s.AuthRepository.FilterStaticAccessTokens(
		ctx,
		authrp.FilterStaticAccessTokenData{
			ID: &[]string{data.ID},
		},
	); err != nil {
		return "", fmt.Errorf("filter: %w", err)
	} else if len(tokens) == 1 {
		return "", tokenssp.ErrStaticTokenExist
	}

	// Create token
	token, err := s.AuthRepository.CreateStaticAccessToken(
		ctx,
		authrp.CreateStaticAccessTokenData(data),
	)
	if err != nil {
		return "", fmt.Errorf("create: %w", err)
	}

	return token, nil
}

// FilterStaticAccessTokens filters static access tokens.
func (s *service) FilterStaticAccessTokens(
	ctx context.Context,
	data tokenssp.FilterStaticAccessTokenData,
) ([]tokenssp.StaticAccessTokenResult, error) {
	// Get tokens
	tokens, err := s.AuthRepository.FilterStaticAccessTokens(
		ctx,
		authrp.FilterStaticAccessTokenData(data),
	)
	if err != nil {
		return nil, fmt.Errorf("filter: %w", err)
	}

	// Map repository to service results
	res := make([]tokenssp.StaticAccessTokenResult, len(tokens))
	for i := range tokens {
		res[i] = tokenssp.StaticAccessTokenResult(tokens[i])
	}

	return res, nil
}

// DeleteStaticAccessToken deletes a static access token by ID.
func (s *service) DeleteStaticAccessToken(
	ctx context.Context,
	id string,
) error {
	// Delete token
	if err := s.AuthRepository.DeleteStaticAccessToken(ctx, id); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (s *service) validateTTL(ttl time.Time) error {
	if time.Now().UTC().After(ttl.Add(s.Config.Auth.MaxClockSkew)) {
		return transport.ErrBadRequest
	}

	return nil
}

// upsertSession creates a new session or updates an existing one.
func (s *service) upsertSession(
	ctx context.Context,
	data *tokenssp.AuthData,
	refreshTokenID string,
) (bool, error) {
	session, err := s.AuthRepository.GetSession(ctx, data.User, data.Device)
	if err != nil && !errors.Is(err, authrp.ErrSessionNotFound) {
		return false, fmt.Errorf("get session: %w", err)
	}

	if session == nil {
		if err := s.AuthRepository.NewSession(
			ctx,
			data.User,
			data.Device,
			&authrp.Session{
				Jti:            refreshTokenID,
				IssuedAt:       time.Now().Format(time.RFC3339),
				Location:       data.MetaLocation,
				IP:             data.MetaIP,
				UserAgent:      data.MetaUserAgent,
				OsFullName:     data.MetaOsFullName,
				OsName:         data.MetaOsName,
				OsVersion:      data.MetaOsVersion,
				Platform:       data.MetaPlatform,
				Model:          data.MetaModel,
				BrowserName:    data.MetaBrowserName,
				BrowserVersion: data.MetaBrowserVersion,
				EngineName:     data.MetaEngineName,
				EngineVersion:  data.MetaEngineVersion,
			},
		); err != nil {
			return false, fmt.Errorf("new session: %w", err)
		}

		return true, nil
	}

	if err := s.AuthRepository.UpdateSession(
		ctx,
		data.User,
		data.Device,
		refreshTokenID,
	); err != nil {
		return false, fmt.Errorf("update session: %w", err)
	}

	return false, nil
}
