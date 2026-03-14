// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package devices

import (
	"context"
	"fmt"
	"log/slog"

	authrp "go.microcore.dev/auth-service/internal/port/adapter/repository/auth"
	devicessp "go.microcore.dev/auth-service/internal/port/service/devices"
)

type (
	// ServiceConfig provides devices service configuration.
	ServiceConfig struct {
		Logger         *slog.Logger
		AuthRepository authrp.Adapter
	}

	service struct {
		*ServiceConfig
	}
)

// NewService creates a new instance of the service.
func NewService(config *ServiceConfig) devicessp.Service {
	return &service{config}
}

// GetActiveDevices returns a list of active devices for a given user.
func (s *service) GetActiveDevices(
	ctx context.Context,
	user uint,
) ([]devicessp.DeviceResult, error) {
	// Get active devices
	devices, err := s.AuthRepository.GetActiveDevices(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	// Make res
	res := make([]devicessp.DeviceResult, len(devices))
	for i := range devices {
		res[i] = devicessp.DeviceResult{
			ID: devices[i].ID,
			Session: devicessp.SessionResult{
				IssuedAt:       devices[i].Session.IssuedAt,
				Location:       devices[i].Session.Location,
				IP:             devices[i].Session.IP,
				UserAgent:      devices[i].Session.UserAgent,
				OsFullName:     devices[i].Session.OsFullName,
				OsName:         devices[i].Session.OsName,
				OsVersion:      devices[i].Session.OsVersion,
				Platform:       devices[i].Session.Platform,
				Model:          devices[i].Session.Model,
				BrowserName:    devices[i].Session.BrowserName,
				BrowserVersion: devices[i].Session.BrowserVersion,
				EngineName:     devices[i].Session.EngineName,
				EngineVersion:  devices[i].Session.EngineVersion,
			},
		}
	}

	return res, nil
}
