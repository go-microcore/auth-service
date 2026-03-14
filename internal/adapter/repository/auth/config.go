// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package auth

import (
	"errors"
	"fmt"
	"time"

	"go.microcore.dev/auth-service/internal/port/adapter/repository/auth"
	"go.microcore.dev/framework/config/env"
	"go.microcore.dev/framework/shutdown"
)

const (
	// DefaultJWTAccessTTL defines the default TTL for access JWT tokens.
	DefaultJWTAccessTTL = 15 * time.Minute

	// DefaultJWTRefreshTTL defines the default TTL for refresh JWT tokens.
	DefaultJWTRefreshTTL = 360 * time.Hour

	// DefaultJWTIssuer specifies the default issuer for JWT tokens.
	DefaultJWTIssuer = "Microcore"

	// DefaultCacheStaticTokenTTL defines how long static tokens are cached.
	DefaultCacheStaticTokenTTL = 1 * time.Minute

	// DefaultCacheLocalTokenSize specifies the default size for local token cache.
	DefaultCacheLocalTokenSize = int(1000)
)

var (
	// ErrJWTAccessKeyTooShort signals that the access JWT key is shorter than
	// the required minimum.
	ErrJWTAccessKeyTooShort = errors.New("jwt access key too short")

	// ErrJWTRefreshKeyTooShort signals that the refresh JWT key is shorter than
	// the required minimum.
	ErrJWTRefreshKeyTooShort = errors.New("jwt refresh key too short")

	// ErrJWTHashKeyTooShort signals that the JWT hash key is shorter than the
	// required minimum.
	ErrJWTHashKeyTooShort = errors.New("jwt hash key too short")

	// ErrAuthKeyTooShort signals that the provided auth key is shorter than the minimum
	// required length.
	ErrAuthKeyTooShort = errors.New("auth key too short")
)

type (
	// Config defines auth configuration.
	Config struct {
		JWT   *ConfigJWT
		Cache *ConfigCache
		Auth  *ConfigAuth
	}

	// ConfigJWT defines jwt configuration.
	ConfigJWT struct {
		AccessKey  []byte        // JWT access token key
		RefreshKey []byte        // JWT refresh token key
		HashKey    []byte        // Key to hash tokens
		AccessTTL  time.Duration // Access token expiration
		RefreshTTL time.Duration // Refresh token expiration
		Issuer     string        // JWT issuer
	}

	// ConfigCache defines cache configuration.
	ConfigCache struct {
		StaticTokenTTL time.Duration // TTL for static token cache
		LocalTokenSize int           // Max items in local token cache
	}

	// ConfigAuth contains authentication configuration.
	ConfigAuth struct {
		Key []byte
	}
)

// NewConfig creates and validates auth configuration.
//
//nolint:funlen // intentionally long: sets all default HTTP server settings
func NewConfig() (*Config, error) {
	var err error

	config := &Config{
		JWT: &ConfigJWT{
			AccessKey:  []byte{},
			RefreshKey: []byte{},
			HashKey:    []byte{},
			AccessTTL: env.DurDefault(
				"JWT_ACCESS_TTL", DefaultJWTAccessTTL,
			),
			RefreshTTL: env.DurDefault(
				"JWT_REFRESH_TTL", DefaultJWTRefreshTTL,
			),
			Issuer: env.StrDefault(
				"JWT_ISSUER", DefaultJWTIssuer,
			),
		},
		Cache: &ConfigCache{
			StaticTokenTTL: env.DurDefault(
				"CACHE_STATIC_TOKEN_TTL", DefaultCacheStaticTokenTTL,
			),
			LocalTokenSize: env.IntDefault(
				"CACHE_LOCAL_TOKEN_SIZE", DefaultCacheLocalTokenSize,
			),
		},
		Auth: &ConfigAuth{
			Key: []byte{},
		},
	}

	// JWT

	if config.JWT.AccessKey, err = env.BytesB64("JWT_ACCESS_KEY"); err != nil {
		return nil, newConfigError(err)
	}

	if len(config.JWT.AccessKey) < auth.JWTSignKeyMinLen {
		return nil, newConfigError(
			fmt.Errorf(
				"%w: must be at least %d bytes",
				ErrJWTAccessKeyTooShort,
				auth.JWTSignKeyMinLen,
			),
		)
	}

	if config.JWT.RefreshKey, err = env.BytesB64("JWT_REFRESH_KEY"); err != nil {
		return nil, newConfigError(err)
	}

	if len(config.JWT.RefreshKey) < auth.JWTSignKeyMinLen {
		return nil, newConfigError(
			fmt.Errorf(
				"%w: must be at least %d bytes",
				ErrJWTRefreshKeyTooShort,
				auth.JWTSignKeyMinLen,
			),
		)
	}

	if config.JWT.HashKey, err = env.BytesB64("JWT_HASH_KEY"); err != nil {
		return nil, newConfigError(err)
	}

	if len(config.JWT.HashKey) < auth.JWTHashKeyMinLen {
		return nil, newConfigError(
			fmt.Errorf(
				"%w: must be at least %d bytes",
				ErrJWTHashKeyTooShort,
				auth.JWTHashKeyMinLen,
			),
		)
	}

	// Auth

	if config.Auth.Key, err = env.BytesB64("AUTH_KEY"); err != nil {
		return nil, newConfigError(err)
	}

	if len(config.Auth.Key) < auth.AuthKeyMinLen {
		return nil, newConfigError(
			fmt.Errorf(
				"%w: must be at least %d bytes",
				ErrAuthKeyTooShort,
				auth.AuthKeyMinLen,
			),
		)
	}

	return config, nil
}

// newConfigError wraps a configuration error as a shutdown exit reason.
func newConfigError(err error) error {
	return shutdown.NewExitReason(
		shutdown.ExitConfigError,
		err,
	)
}
