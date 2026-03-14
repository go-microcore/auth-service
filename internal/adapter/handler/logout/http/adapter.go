// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package http

import (
	"context"
	"fmt"
	"log/slog"

	logouthp "go.microcore.dev/auth-service/internal/port/adapter/handler/logout/http"
	logoutsp "go.microcore.dev/auth-service/internal/port/service/logout"
	"go.microcore.dev/framework/transport/http"
	"go.microcore.dev/framework/transport/http/server"
)

type (
	// AdapterConfig provides logout HTTP adapter handler configuration.
	AdapterConfig struct {
		Logger        *slog.Logger
		LogoutService logoutsp.Service
	}

	adapter struct {
		*AdapterConfig
	}
)

// NewAdapter creates a new instance of the adapter.
func NewAdapter(config *AdapterConfig) logouthp.Adapter {
	return &adapter{config}
}

// Logout terminate the current session on user device.
//
// @Summary Terminate the current session on user device.
// @Tags Logout
// @Security BearerAuth
// @Produce json
// @Success 204
// @Failure 401 {object} server.ErrResponse "Possible codes: INVALID_TOKEN, 2FA_REQUIRED"
// @Failure 403 {object} server.ErrResponse "Possible codes: INSUFFICIENT_PERMISSIONS"
// @Failure 503 {object} server.ErrResponse "Possible codes: SERVICE_UNAVAILABLE"
// @Router /auth/logout/ [post]
func (a *adapter) Logout(
	ctx context.Context,
	c *server.RequestContext,
) {
	// Get user
	user, err := c.UserValueUint("user")
	if err != nil {
		c.WriteError(fmt.Errorf("get user: %w", err))
		return
	}

	// Get device
	device, err := c.UserValueStr("device")
	if err != nil {
		c.WriteError(fmt.Errorf("get device: %w", err))
		return
	}

	// Logout current device
	if err := a.LogoutService.LogoutDevice(
		ctx,
		logoutsp.DeviceData{
			User:   user,
			Device: device,
		},
	); err != nil {
		c.WriteError(fmt.Errorf("logout device: %w", err))
		return
	}

	// Write res
	c.StatusCode(http.StatusNoContent)
}

// LogoutAll terminate all active sessions for a user across all devices.
//
// @Summary Terminate all active sessions for a user across all devices.
// @Tags Logout
// @Security BearerAuth
// @Produce json
// @Success 204
// @Failure 401 {object} server.ErrResponse "Possible codes: INVALID_TOKEN, 2FA_REQUIRED"
// @Failure 403 {object} server.ErrResponse "Possible codes: INSUFFICIENT_PERMISSIONS"
// @Failure 503 {object} server.ErrResponse "Possible codes: SERVICE_UNAVAILABLE"
// @Router /auth/logout/all [post]
func (a *adapter) LogoutAll(
	ctx context.Context,
	c *server.RequestContext,
) {
	// Get user
	user, err := c.UserValueUint("user")
	if err != nil {
		c.WriteError(fmt.Errorf("get user: %w", err))
		return
	}

	// Logout all devices
	if err := a.LogoutService.LogoutAll(
		ctx,
		logoutsp.AllData{
			User: user,
		},
	); err != nil {
		c.WriteError(fmt.Errorf("logout all: %w", err))
		return
	}

	// Write res
	c.StatusCode(http.StatusNoContent)
}

// LogoutDevice terminate the session on a selected user device.
//
// @Summary Terminate the session on a selected user device.
// @Tags Logout
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body logouthp.LogoutDeviceRequest true "Request data"
// @Success 204
// @Failure 400 {object} server.ErrResponse "Possible codes: INVALID_DEVICE"
// @Failure 401 {object} server.ErrResponse "Possible codes: INVALID_TOKEN, 2FA_REQUIRED"
// @Failure 403 {object} server.ErrResponse "Possible codes: INSUFFICIENT_PERMISSIONS"
// @Failure 415 {object} server.ErrResponse "Possible codes: INVALID_JSON_BODY"
// @Failure 503 {object} server.ErrResponse "Possible codes: SERVICE_UNAVAILABLE"
// @Router /auth/logout/device [post]
func (a *adapter) LogoutDevice(
	ctx context.Context,
	c *server.RequestContext,
	req *logouthp.LogoutDeviceRequest,
) {
	// Get user
	user, err := c.UserValueUint("user")
	if err != nil {
		c.WriteError(fmt.Errorf("get user: %w", err))
		return
	}

	// Logout target device
	if err := a.LogoutService.LogoutDevice(
		ctx,
		logoutsp.DeviceData{
			User:   user,
			Device: req.Device,
		},
	); err != nil {
		c.WriteError(fmt.Errorf("logout device: %w", err))
		return
	}

	// Write res
	c.StatusCode(http.StatusNoContent)
}
