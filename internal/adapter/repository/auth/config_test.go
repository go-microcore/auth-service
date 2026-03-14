// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package auth_test

import (
	"encoding/base64"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.microcore.dev/auth-service/internal/adapter/repository/auth"
)

func TestNewConfig(t *testing.T) {
	// JWT
	t.Setenv("JWT_ACCESS_KEY", "JF3fzHUr/QIT4ylO5avPGhu6FJE/SK1AnZhODfWJUcH8EdPq+hbYvTQ9cy1QBqzZRrFjaKeirk/VHf/ZB4QkaA==")
	t.Setenv("JWT_REFRESH_KEY", "7KfglXcOWsioJDdeLsrZ21YC4vAYp0PW1G3SLtZKSfZd2BXs1QULgf+gpwDZ5SgdobsN206LknyDTsr+WcIjMQ==")
	t.Setenv("JWT_HASH_KEY", "UMFqRKnRiUkPUlqtok4TWc1YFJPtrJOGno8YPX5C1Dc=")
	t.Setenv("JWT_ACCESS_TTL", "1m")
	t.Setenv("JWT_REFRESH_TTL", "1h")
	t.Setenv("JWT_ISSUER", "test")
	// Cache
	t.Setenv("CACHE_STATIC_TOKEN_TTL", "2m")
	t.Setenv("CACHE_LOCAL_TOKEN_SIZE", "100")
	// Auth
	t.Setenv("AUTH_KEY", "HwaATOL8yS+2hAz1w8jR8s0DwSZ00k/NvZy5u1s7t9c=")

	jwtAccessKey, err := base64.StdEncoding.DecodeString(os.Getenv("JWT_ACCESS_KEY"))
	require.NoError(t, err)

	jwtRefreshKey, err := base64.StdEncoding.DecodeString(os.Getenv("JWT_REFRESH_KEY"))
	require.NoError(t, err)

	jwtHashKey, err := base64.StdEncoding.DecodeString(os.Getenv("JWT_HASH_KEY"))
	require.NoError(t, err)

	jwtAccessTTL, err := time.ParseDuration(os.Getenv("JWT_ACCESS_TTL"))
	require.NoError(t, err)

	jwtRefreshTTL, err := time.ParseDuration(os.Getenv("JWT_REFRESH_TTL"))
	require.NoError(t, err)

	cacheStaticTokenTTL, err := time.ParseDuration(os.Getenv("CACHE_STATIC_TOKEN_TTL"))
	require.NoError(t, err)

	cacheLocalTokenSize, err := strconv.Atoi(os.Getenv("CACHE_LOCAL_TOKEN_SIZE"))
	require.NoError(t, err)

	authKey, err := base64.StdEncoding.DecodeString(os.Getenv("AUTH_KEY"))
	require.NoError(t, err)

	cfg, err := auth.NewConfig()
	require.NoError(t, err)

	// JWT
	require.Equal(t, jwtAccessKey, cfg.JWT.AccessKey)
	require.Equal(t, jwtRefreshKey, cfg.JWT.RefreshKey)
	require.Equal(t, jwtHashKey, cfg.JWT.HashKey)
	require.Equal(t, jwtAccessTTL, cfg.JWT.AccessTTL)
	require.Equal(t, jwtRefreshTTL, cfg.JWT.RefreshTTL)
	require.Equal(t, "test", cfg.JWT.Issuer)
	// Cache
	require.Equal(t, cacheStaticTokenTTL, cfg.Cache.StaticTokenTTL)
	require.Equal(t, cacheLocalTokenSize, cfg.Cache.LocalTokenSize)
	// Auth
	require.Equal(t, authKey, cfg.Auth.Key)
}

func TestNewConfig_Default(t *testing.T) {
	t.Setenv("JWT_ACCESS_KEY", "JF3fzHUr/QIT4ylO5avPGhu6FJE/SK1AnZhODfWJUcH8EdPq+hbYvTQ9cy1QBqzZRrFjaKeirk/VHf/ZB4QkaA==")
	t.Setenv("JWT_REFRESH_KEY", "7KfglXcOWsioJDdeLsrZ21YC4vAYp0PW1G3SLtZKSfZd2BXs1QULgf+gpwDZ5SgdobsN206LknyDTsr+WcIjMQ==")
	t.Setenv("JWT_HASH_KEY", "UMFqRKnRiUkPUlqtok4TWc1YFJPtrJOGno8YPX5C1Dc=")
	t.Setenv("AUTH_KEY", "HwaATOL8yS+2hAz1w8jR8s0DwSZ00k/NvZy5u1s7t9c=")

	cfg, err := auth.NewConfig()
	require.NoError(t, err)

	// JWT
	require.Equal(t, auth.DefaultJWTAccessTTL, cfg.JWT.AccessTTL)
	require.Equal(t, auth.DefaultJWTRefreshTTL, cfg.JWT.RefreshTTL)
	require.Equal(t, auth.DefaultJWTIssuer, cfg.JWT.Issuer)
	// Cache
	require.Equal(t, auth.DefaultCacheStaticTokenTTL, cfg.Cache.StaticTokenTTL)
	require.Equal(t, auth.DefaultCacheLocalTokenSize, cfg.Cache.LocalTokenSize)
}

func TestNewConfig_Error(t *testing.T) {
	t.Setenv("JWT_ACCESS_KEY", "JF3fzHUr/QIT4ylO5avPGhu6FJE/SK1AnZhODfWJUcH8EdPq+hbYvTQ9cy1QBqzZRrFjaKeirk/VHf/ZB4QkaA==")
	t.Setenv("JWT_REFRESH_KEY", "7KfglXcOWsioJDdeLsrZ21YC4vAYp0PW1G3SLtZKSfZd2BXs1QULgf+gpwDZ5SgdobsN206LknyDTsr+WcIjMQ==")
	t.Setenv("JWT_HASH_KEY", "UMFqRKnRiUkPUlqtok4TWc1YFJPtrJOGno8YPX5C1Dc=")
	t.Setenv("AUTH_KEY", "HwaATOL8yS+2hAz1w8jR8s0DwSZ00k/NvZy5u1s7t9c=")

	_, err := auth.NewConfig()
	require.NoError(t, err)

	t.Run("JWT_ACCESS_KEY no pass", func(t *testing.T) {
		t.Setenv("JWT_ACCESS_KEY", "")

		_, err := auth.NewConfig()
		require.Error(t, err)
		require.ErrorContains(t, err, "variable JWT_ACCESS_KEY is not set")
	})

	t.Run("JWT_REFRESH_KEY no pass", func(t *testing.T) {
		t.Setenv("JWT_REFRESH_KEY", "")

		_, err := auth.NewConfig()
		require.Error(t, err)
		require.ErrorContains(t, err, "variable JWT_REFRESH_KEY is not set")
	})

	t.Run("JWT_HASH_KEY no pass", func(t *testing.T) {
		t.Setenv("JWT_HASH_KEY", "")

		_, err := auth.NewConfig()
		require.Error(t, err)
		require.ErrorContains(t, err, "variable JWT_HASH_KEY is not set")
	})

	t.Run("AUTH_KEY no pass", func(t *testing.T) {
		t.Setenv("AUTH_KEY", "")

		_, err := auth.NewConfig()
		require.Error(t, err)
		require.ErrorContains(t, err, "variable AUTH_KEY is not set")
	})
}
