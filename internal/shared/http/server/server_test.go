// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package server_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	sharedhttp "go.microcore.dev/auth-service/internal/shared/http/server"
	sharedtel "go.microcore.dev/auth-service/internal/shared/telemetry"
	"go.microcore.dev/framework/transport/http/server"
)

func TestSetup_UseCors(t *testing.T) {
	t.Parallel()

	mockHTTPServerManager := server.NewMockManager(t)

	mockHTTPServerManager.EXPECT().
		UseCors(mock.Anything).
		Return(nil)

	httpServer := &sharedhttp.Server{
		Config: &sharedhttp.Config{
			Name:               "",
			Concurrency:        0,
			ReadBufferSize:     0,
			WriteBufferSize:    0,
			ReadTimeout:        0,
			WriteTimeout:       0,
			IdleTimeout:        0,
			MaxConnsPerIP:      0,
			MaxRequestsPerConn: 0,
			MaxRequestBodySize: 0,
			DisableKeepalive:   false,
			TCPKeepalive:       false,
			LogAllErrors:       false,
			Cors: &sharedhttp.ConfigCors{
				Origin:  "",
				Methods: "",
				Headers: "",
			},
			Host:    "",
			Port:    "",
			Swagger: false,
		},
		Telemetry: &sharedtel.Telemetry{
			Options: &sharedtel.Options{
				AppName: "",
			},
			Config: &sharedtel.Config{
				Enabled:  false,
				Endpoint: "",
			},
			Manager: nil,
		},
		Manager: mockHTTPServerManager,
	}

	httpServer.Setup()
}

func TestSetup_WithSwagger(t *testing.T) {
	t.Parallel()

	mockHTTPServerManager := server.NewMockManager(t)

	mockHTTPServerManager.EXPECT().
		UseCors(mock.Anything).
		Return(nil)

	mockHTTPServerManager.EXPECT().
		UseSwagger().
		Return(nil)

	httpServer := &sharedhttp.Server{
		Config: &sharedhttp.Config{
			Name:               "",
			Concurrency:        0,
			ReadBufferSize:     0,
			WriteBufferSize:    0,
			ReadTimeout:        0,
			WriteTimeout:       0,
			IdleTimeout:        0,
			MaxConnsPerIP:      0,
			MaxRequestsPerConn: 0,
			MaxRequestBodySize: 0,
			DisableKeepalive:   false,
			TCPKeepalive:       false,
			LogAllErrors:       false,
			Cors: &sharedhttp.ConfigCors{
				Origin:  "",
				Methods: "",
				Headers: "",
			},
			Host:    "",
			Port:    "",
			Swagger: true,
		},
		Telemetry: &sharedtel.Telemetry{
			Options: &sharedtel.Options{
				AppName: "",
			},
			Config: &sharedtel.Config{
				Enabled:  false,
				Endpoint: "",
			},
			Manager: nil,
		},
		Manager: mockHTTPServerManager,
	}

	httpServer.Setup()
}

func TestSetup_WithTelemetry(t *testing.T) {
	t.Parallel()

	mockHTTPServerManager := server.NewMockManager(t)

	mockHTTPServerManager.EXPECT().
		UseCors(mock.Anything).
		Return(nil)

	mockHTTPServerManager.EXPECT().
		SetTelemetryManager(mock.Anything).
		Return(nil)

	httpServer := &sharedhttp.Server{
		Config: &sharedhttp.Config{
			Name:               "",
			Concurrency:        0,
			ReadBufferSize:     0,
			WriteBufferSize:    0,
			ReadTimeout:        0,
			WriteTimeout:       0,
			IdleTimeout:        0,
			MaxConnsPerIP:      0,
			MaxRequestsPerConn: 0,
			MaxRequestBodySize: 0,
			DisableKeepalive:   false,
			TCPKeepalive:       false,
			LogAllErrors:       false,
			Cors: &sharedhttp.ConfigCors{
				Origin:  "",
				Methods: "",
				Headers: "",
			},
			Host:    "",
			Port:    "",
			Swagger: false,
		},
		Telemetry: &sharedtel.Telemetry{
			Options: &sharedtel.Options{
				AppName: "",
			},
			Config: &sharedtel.Config{
				Enabled:  true,
				Endpoint: "",
			},
			Manager: nil,
		},
		Manager: mockHTTPServerManager,
	}

	httpServer.Setup()
}
