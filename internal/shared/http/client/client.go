// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package client

import (
	"go.microcore.dev/auth-service/internal/shared/telemetry"
	"go.microcore.dev/framework/transport/http/client"
)

type (
	// Client represents a HTTP client helper object.
	Client struct {
		Telemetry *telemetry.Telemetry
		Manager   client.Manager
	}
)

// Setup configures the helper object.
func (c *Client) Setup() {
	// Set telemetry manager
	if c.Telemetry.Config.Enabled {
		c.Manager.SetTelemetryManager(c.Telemetry.Manager)
	}
}
