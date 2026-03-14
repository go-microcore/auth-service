// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package postgres_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	sharedpg "go.microcore.dev/auth-service/internal/shared/postgres"
	sharedtel "go.microcore.dev/auth-service/internal/shared/telemetry"
	"go.microcore.dev/framework/db/postgres"
)

func TestSetup_WithoutTelemetry(t *testing.T) {
	t.Parallel()

	r := &sharedpg.Postgres{
		Config: &sharedpg.Config{
			Host:     "",
			Port:     "",
			User:     "",
			Password: "",
			DB:       "",
			SSL:      "",
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
		Manager: nil,
	}

	err := r.Setup()
	require.NoError(t, err)
}

func TestSetup_WithTelemetry(t *testing.T) {
	t.Parallel()

	mockPostgresManager := postgres.NewMockManager(t)

	mockPostgresManager.EXPECT().
		SetTelemetryManager(mock.Anything).
		Return(nil)

	p := &sharedpg.Postgres{
		Config: &sharedpg.Config{
			Host:     "",
			Port:     "",
			User:     "",
			Password: "",
			DB:       "",
			SSL:      "",
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
		Manager: mockPostgresManager,
	}

	err := p.Setup()
	require.NoError(t, err)
}
