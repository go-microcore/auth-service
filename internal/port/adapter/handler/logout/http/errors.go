// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package http

import (
	"go.microcore.dev/framework/transport"
)

// ErrInvalidDevice is returned when a device identifier in a request is
// invalid or missing.
var ErrInvalidDevice = transport.NewError(
	transport.ErrBadRequest,
	"invalid device",
	"INVALID_DEVICE",
)
