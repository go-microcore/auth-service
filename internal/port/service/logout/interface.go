// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package logout

import (
	"context"
)

type (
	// Service defines the interface for the logout service.
	Service interface {
		LogoutDevice(
			ctx context.Context,
			data DeviceData,
		) error

		LogoutAll(
			ctx context.Context,
			data AllData,
		) error
	}
)
