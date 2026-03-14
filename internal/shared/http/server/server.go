// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package server

import (
	"go.microcore.dev/auth-service/internal/shared/telemetry"
	"go.microcore.dev/framework/transport/http/server"
)

type (
	// Server represents a HTTP server helper object.
	Server struct {
		Config    *Config
		Telemetry *telemetry.Telemetry
		Manager   server.Manager
	}
)

// Setup configures the helper object.
func (s *Server) Setup() {
	// Use CORS
	s.Manager.UseCors(
		server.WithCorsOrigin(s.Config.Cors.Origin),
		server.WithCorsMethods(s.Config.Cors.Methods),
		server.WithCorsHeaders(s.Config.Cors.Headers),
	)

	// Use Swagger
	if s.Config.Swagger {
		s.Manager.UseSwagger()
	}

	// Set telemetry manager
	if s.Telemetry.Config.Enabled {
		s.Manager.SetTelemetryManager(s.Telemetry.Manager)
	}
}
