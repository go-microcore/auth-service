// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package http

type (
	// Requests

	// LogoutDeviceRequest represents the API request to log out a specific device.
	LogoutDeviceRequest struct {
		Device string `json:"device"`
	}
)

// Validate checks that the request contains a valid device identifier.
func (r *LogoutDeviceRequest) Validate() error {
	if r.Device == "" {
		return ErrInvalidDevice
	}

	return nil
}
