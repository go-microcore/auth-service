// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package telemetry

import (
	"go.microcore.dev/framework/telemetry"
)

type (
	// Telemetry represents a Telemetry helper object.
	Telemetry struct {
		Options *Options
		Config  *Config
		Manager telemetry.Manager
	}

	// Options defines options for creating a Telemetry helper object.
	Options struct {
		AppName string
	}
)
