// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package redis_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	sharedredis "go.microcore.dev/auth-service/internal/shared/redis"
	sharedtel "go.microcore.dev/auth-service/internal/shared/telemetry"
	"go.microcore.dev/framework/db/redis"
)

func TestSetup_WithoutTelemetry(t *testing.T) {
	t.Parallel()

	r := &sharedredis.Redis{
		Config: &sharedredis.Config{
			Addr:     "",
			Password: "",
			DB:       0,
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

	mockRedisManager := redis.NewMockManager(t)

	mockRedisManager.EXPECT().
		SetTelemetryManager(mock.Anything).
		Return(nil)

	r := &sharedredis.Redis{
		Config: &sharedredis.Config{
			Addr:     "",
			Password: "",
			DB:       0,
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
		Manager: mockRedisManager,
	}

	err := r.Setup()
	require.NoError(t, err)
}
