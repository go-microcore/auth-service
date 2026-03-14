// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package logout

import (
	"context"
	"fmt"
	"log/slog"

	authrp "go.microcore.dev/auth-service/internal/port/adapter/repository/auth"
	logoutsp "go.microcore.dev/auth-service/internal/port/service/logout"
)

type (
	// ServiceConfig provides logout service configuration.
	ServiceConfig struct {
		Logger         *slog.Logger
		AuthRepository authrp.Adapter
	}

	service struct {
		*ServiceConfig
	}
)

// NewService creates a new instance of the service.
func NewService(config *ServiceConfig) logoutsp.Service {
	return &service{config}
}

// LogoutDevice ends the session on a selected user device.
func (s *service) LogoutDevice(
	ctx context.Context,
	data logoutsp.DeviceData,
) error {
	if err := s.AuthRepository.DeleteSession(
		ctx,
		data.User,
		data.Device,
	); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}

// LogoutAll ends all active sessions for a user across all devices.
func (s *service) LogoutAll(
	ctx context.Context,
	data logoutsp.AllData,
) error {
	// Get active devices
	devices, err := s.AuthRepository.GetActiveDevices(ctx, data.User)
	if err != nil {
		return fmt.Errorf("get active devices: %w", err)
	}

	// Delete active devices
	for i := range devices {
		if err := s.AuthRepository.DeleteSession(
			ctx,
			data.User,
			devices[i].ID,
		); err != nil {
			return fmt.Errorf("delete session: %w", err)
		}
	}

	return nil
}
