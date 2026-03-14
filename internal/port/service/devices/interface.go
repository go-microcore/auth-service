// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package devices

import (
	"context"
)

type (
	// Service defines the interface for the devices service.
	Service interface {
		GetActiveDevices(
			ctx context.Context,
			user uint,
		) ([]DeviceResult, error)
	}
)
