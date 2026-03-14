// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package client_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	sharedhttp "go.microcore.dev/auth-service/internal/shared/http/client"
	sharedtel "go.microcore.dev/auth-service/internal/shared/telemetry"
	"go.microcore.dev/framework/transport/http/client"
)

func TestSetup_WithoutTelemetry(t *testing.T) {
	t.Parallel()

	httpClient := &sharedhttp.Client{
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
		Manager: client.NewMockManager(t),
	}

	httpClient.Setup()
}

func TestSetup_WithTelemetry(t *testing.T) {
	t.Parallel()

	mockHTTPClientManager := client.NewMockManager(t)

	mockHTTPClientManager.EXPECT().
		SetTelemetryManager(mock.Anything).
		Return(nil)

	httpClient := &sharedhttp.Client{
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
		Manager: mockHTTPClientManager,
	}

	httpClient.Setup()
}
