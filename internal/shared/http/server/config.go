// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package server

import (
	"time"

	"go.microcore.dev/framework/config/env"
)

const (
	// DefaultHTTPServerName is the default name of the HTTP server.
	DefaultHTTPServerName = "microcore"

	// DefaultHTTPServerConcurrency is the default max number of concurrent requests.
	DefaultHTTPServerConcurrency = 1024

	// DefaultHTTPServerReadBufferSize is the default read buffer size in bytes.
	DefaultHTTPServerReadBufferSize = 1024 * 4 // 4KB

	// DefaultHTTPServerWriteBufferSize is the default write buffer size in bytes.
	DefaultHTTPServerWriteBufferSize = 1024 * 4 // 4KB

	// DefaultHTTPServerReadTimeout is the default maximum duration for reading the full request.
	DefaultHTTPServerReadTimeout = 5 * time.Second

	// DefaultHTTPServerWriteTimeout is the default maximum duration before timing out writes.
	DefaultHTTPServerWriteTimeout = 5 * time.Second

	// DefaultHTTPServerIdleTimeout is the default maximum idle time for connections.
	DefaultHTTPServerIdleTimeout = 10 * time.Second

	// DefaultHTTPServerMaxConnsPerIP is the default maximum connections per IP.
	DefaultHTTPServerMaxConnsPerIP = 0

	// DefaultHTTPServerMaxRequestsPerConn is the default maximum requests per connection.
	DefaultHTTPServerMaxRequestsPerConn = 0

	// DefaultHTTPServerMaxRequestBodySize is the default maximum request body size in bytes.
	DefaultHTTPServerMaxRequestBodySize = 1024 * 1024 * 8 // 8MB

	// DefaultHTTPServerDisableKeepalive specifies if keep-alives are disabled by default.
	DefaultHTTPServerDisableKeepalive = false

	// DefaultHTTPServerTCPKeepalive specifies if TCP keep-alive is enabled by default.
	DefaultHTTPServerTCPKeepalive = false

	// DefaultHTTPServerLogAllErrors specifies if all server errors are logged by default.
	DefaultHTTPServerLogAllErrors = false

	// DefaultHTTPServerCorsOrigin is the default CORS allowed origin.
	DefaultHTTPServerCorsOrigin = "*"

	// DefaultHTTPServerCorsMethods is the default CORS allowed methods.
	DefaultHTTPServerCorsMethods = "*"

	// DefaultHTTPServerCorsHeaders is the default CORS allowed headers.
	DefaultHTTPServerCorsHeaders = "*"

	// DefaultHTTPServerHost is the default host to bind the server.
	DefaultHTTPServerHost = "0.0.0.0"

	// DefaultHTTPServerPort is the default port to bind the server.
	DefaultHTTPServerPort = "80"

	// DefaultHTTPServerSwagger specifies if Swagger UI is enabled by default.
	DefaultHTTPServerSwagger = true
)

type (
	// Config defines the HTTP server configuration.
	Config struct {
		Name               string
		Concurrency        int
		ReadBufferSize     int
		WriteBufferSize    int
		ReadTimeout        time.Duration
		WriteTimeout       time.Duration
		IdleTimeout        time.Duration
		MaxConnsPerIP      int
		MaxRequestsPerConn int
		MaxRequestBodySize int
		DisableKeepalive   bool
		TCPKeepalive       bool
		LogAllErrors       bool
		Cors               *ConfigCors
		Host               string
		Port               string
		Swagger            bool
	}

	// ConfigCors defines the CORS configuration.
	ConfigCors struct {
		Origin  string
		Methods string
		Headers string
	}
)

// NewConfig creates and validates a HTTP server configuration.
//
//nolint:funlen // intentionally long: sets all default HTTP server settings
func NewConfig() *Config {
	return &Config{
		Name: env.StrDefault(
			"HTTP_SERVER_NAME",
			DefaultHTTPServerName,
		),
		Concurrency: env.IntDefault(
			"HTTP_SERVER_CONCURRENCY",
			DefaultHTTPServerConcurrency,
		),
		ReadBufferSize: env.IntDefault(
			"HTTP_SERVER_READ_BUFFER_SIZE",
			DefaultHTTPServerReadBufferSize,
		),
		WriteBufferSize: env.IntDefault(
			"HTTP_SERVER_WRITE_BUFFER_SIZE",
			DefaultHTTPServerWriteBufferSize,
		),
		ReadTimeout: env.DurDefault(
			"HTTP_SERVER_READ_TIMEOUT",
			DefaultHTTPServerReadTimeout,
		),
		WriteTimeout: env.DurDefault(
			"HTTP_SERVER_WRITE_TIMEOUT",
			DefaultHTTPServerWriteTimeout,
		),
		IdleTimeout: env.DurDefault(
			"HTTP_SERVER_IDLE_TIMEOUT",
			DefaultHTTPServerIdleTimeout,
		),
		MaxConnsPerIP: env.IntDefault(
			"HTTP_SERVER_MAX_CONNS_PER_IP",
			DefaultHTTPServerMaxConnsPerIP,
		),
		MaxRequestsPerConn: env.IntDefault(
			"HTTP_SERVER_MAX_REQUESTS_PER_CONN",
			DefaultHTTPServerMaxRequestsPerConn,
		),
		MaxRequestBodySize: env.IntDefault(
			"HTTP_SERVER_MAX_REQUEST_BODY_SIZE",
			DefaultHTTPServerMaxRequestBodySize,
		),
		DisableKeepalive: env.BoolDefault(
			"HTTP_SERVER_DISABLE_KEEPALIVE",
			DefaultHTTPServerDisableKeepalive,
		),
		TCPKeepalive: env.BoolDefault(
			"HTTP_SERVER_TCP_KEEPALIVE",
			DefaultHTTPServerTCPKeepalive,
		),
		LogAllErrors: env.BoolDefault(
			"HTTP_SERVER_LOG_ALL_ERRORS",
			DefaultHTTPServerLogAllErrors,
		),
		Cors: &ConfigCors{
			Origin: env.StrDefault(
				"HTTP_SERVER_CORS_ORIGIN",
				DefaultHTTPServerCorsOrigin,
			),
			Methods: env.StrDefault(
				"HTTP_SERVER_CORS_METHODS",
				DefaultHTTPServerCorsMethods,
			),
			Headers: env.StrDefault(
				"HTTP_SERVER_CORS_HEADERS",
				DefaultHTTPServerCorsHeaders,
			),
		},
		Host: env.StrDefault(
			"HTTP_SERVER_HOST",
			DefaultHTTPServerHost,
		),
		Port: env.StrDefault(
			"HTTP_SERVER_PORT",
			DefaultHTTPServerPort,
		),
		Swagger: env.BoolDefault(
			"HTTP_SERVER_SWAGGER",
			DefaultHTTPServerSwagger,
		),
	}
}
