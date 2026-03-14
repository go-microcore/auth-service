// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package http

import (
	"context"
	"fmt"
	"log/slog"

	deviceshp "go.microcore.dev/auth-service/internal/port/adapter/handler/devices/http"
	devicessp "go.microcore.dev/auth-service/internal/port/service/devices"
	"go.microcore.dev/framework/transport/http"
	"go.microcore.dev/framework/transport/http/server"
)

type (
	// AdapterConfig provides devices HTTP adapter handler configuration.
	AdapterConfig struct {
		Logger         *slog.Logger
		DevicesService devicessp.Service
	}

	adapter struct {
		*AdapterConfig
	}
)

// NewAdapter creates a new instance of the adapter.
func NewAdapter(config *AdapterConfig) deviceshp.Adapter {
	return &adapter{config}
}

// ListDevices returns a list of user's active devices.
//
// @Summary Returns a list of user's active devices.
// @Tags Devices
// @Security BearerAuth
// @Produce json
// @Success 200 {array} deviceshp.DeviceResponse
// @Failure 401 {object} server.ErrResponse "Possible codes: INVALID_TOKEN, 2FA_REQUIRED"
// @Failure 403 {object} server.ErrResponse "Possible codes: INSUFFICIENT_PERMISSIONS"
// @Failure 503 {object} server.ErrResponse "Possible codes: SERVICE_UNAVAILABLE"
// @Router /auth/devices [get]
func (a *adapter) ListDevices(
	ctx context.Context,
	c *server.RequestContext,
) {
	// Get user
	user, err := c.UserValueUint("user")
	if err != nil {
		c.WriteError(fmt.Errorf("get user: %w", err))
		return
	}

	// Get active devices
	devices, err := a.DevicesService.GetActiveDevices(ctx, user)
	if err != nil {
		c.WriteError(fmt.Errorf("get active devices: %w", err))
		return
	}

	// Make res
	res := make([]deviceshp.DeviceResponse, len(devices))
	for i := range devices {
		res[i] = deviceshp.DeviceResponse{
			ID:      devices[i].ID,
			Session: deviceshp.SessionResponse(devices[i].Session),
		}
	}

	// Write res
	c.WriteJsonWithStatusCode(http.StatusOK, res)
}
