// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package server_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.microcore.dev/auth-service/internal/shared/http/server"
)

func TestNewConfig(t *testing.T) {
	t.Setenv("HTTP_SERVER_NAME", "name")
	t.Setenv("HTTP_SERVER_CONCURRENCY", "100")
	t.Setenv("HTTP_SERVER_READ_BUFFER_SIZE", "1024")
	t.Setenv("HTTP_SERVER_WRITE_BUFFER_SIZE", "1024")
	t.Setenv("HTTP_SERVER_READ_TIMEOUT", "1m")
	t.Setenv("HTTP_SERVER_WRITE_TIMEOUT", "1m")
	t.Setenv("HTTP_SERVER_IDLE_TIMEOUT", "1m")
	t.Setenv("HTTP_SERVER_MAX_CONNS_PER_IP", "100")
	t.Setenv("HTTP_SERVER_MAX_REQUESTS_PER_CONN", "100")
	t.Setenv("HTTP_SERVER_MAX_REQUEST_BODY_SIZE", "1024")
	t.Setenv("HTTP_SERVER_DISABLE_KEEPALIVE", "true")
	t.Setenv("HTTP_SERVER_TCP_KEEPALIVE", "true")
	t.Setenv("HTTP_SERVER_LOG_ALL_ERRORS", "true")
	t.Setenv("HTTP_SERVER_CORS_ORIGIN", "orign")
	t.Setenv("HTTP_SERVER_CORS_METHODS", "methods")
	t.Setenv("HTTP_SERVER_CORS_HEADERS", "headers")
	t.Setenv("HTTP_SERVER_HOST", "127.0.0.1")
	t.Setenv("HTTP_SERVER_PORT", "80")
	t.Setenv("HTTP_SERVER_SWAGGER", "true")

	cfg := server.NewConfig()

	require.Equal(t, "name", cfg.Name)
	require.Equal(t, 100, cfg.Concurrency)
	require.Equal(t, 1024, cfg.ReadBufferSize)
	require.Equal(t, 1024, cfg.WriteBufferSize)
	require.Equal(t, 1*time.Minute, cfg.ReadTimeout)
	require.Equal(t, 1*time.Minute, cfg.WriteTimeout)
	require.Equal(t, 1*time.Minute, cfg.IdleTimeout)
	require.Equal(t, 100, cfg.MaxConnsPerIP)
	require.Equal(t, 100, cfg.MaxRequestsPerConn)
	require.Equal(t, 1024, cfg.MaxRequestBodySize)
	require.True(t, cfg.DisableKeepalive)
	require.True(t, cfg.TCPKeepalive)
	require.True(t, cfg.LogAllErrors)
	require.Equal(t, "orign", cfg.Cors.Origin)
	require.Equal(t, "methods", cfg.Cors.Methods)
	require.Equal(t, "headers", cfg.Cors.Headers)
	require.Equal(t, "127.0.0.1", cfg.Host)
	require.Equal(t, "80", cfg.Port)
	require.True(t, cfg.Swagger)
}

//nolint:paralleltest // uses environment variables
func TestNewConfig_Default(t *testing.T) {
	cfg := server.NewConfig()

	require.Equal(t, server.DefaultHTTPServerName, cfg.Name)
	require.Equal(t, server.DefaultHTTPServerConcurrency, cfg.Concurrency)
	require.Equal(t, server.DefaultHTTPServerReadBufferSize, cfg.ReadBufferSize)
	require.Equal(t, server.DefaultHTTPServerWriteBufferSize, cfg.WriteBufferSize)
	require.Equal(t, server.DefaultHTTPServerReadTimeout, cfg.ReadTimeout)
	require.Equal(t, server.DefaultHTTPServerWriteTimeout, cfg.WriteTimeout)
	require.Equal(t, server.DefaultHTTPServerIdleTimeout, cfg.IdleTimeout)
	require.Equal(t, server.DefaultHTTPServerMaxConnsPerIP, cfg.MaxConnsPerIP)
	require.Equal(t, server.DefaultHTTPServerMaxRequestsPerConn, cfg.MaxRequestsPerConn)
	require.Equal(t, server.DefaultHTTPServerMaxRequestBodySize, cfg.MaxRequestBodySize)
	require.Equal(t, server.DefaultHTTPServerDisableKeepalive, cfg.DisableKeepalive)
	require.Equal(t, server.DefaultHTTPServerTCPKeepalive, cfg.TCPKeepalive)
	require.Equal(t, server.DefaultHTTPServerLogAllErrors, cfg.LogAllErrors)
	require.Equal(t, server.DefaultHTTPServerCorsOrigin, cfg.Cors.Origin)
	require.Equal(t, server.DefaultHTTPServerCorsMethods, cfg.Cors.Methods)
	require.Equal(t, server.DefaultHTTPServerCorsHeaders, cfg.Cors.Headers)
	require.Equal(t, server.DefaultHTTPServerHost, cfg.Host)
	require.Equal(t, server.DefaultHTTPServerPort, cfg.Port)
	require.Equal(t, server.DefaultHTTPServerSwagger, cfg.Swagger)
}
