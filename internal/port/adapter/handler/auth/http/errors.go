// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package http

import (
	"go.microcore.dev/framework/transport"
)

var (
	// ErrInvalidToken is returned when the token is invalid.
	ErrInvalidToken = transport.NewError(
		transport.ErrUnauthorized,
		"invalid token",
		"INVALID_TOKEN",
	)

	// ErrMFARequired is returned when two-factor authentication is required.
	ErrMFARequired = transport.NewError(
		transport.ErrUnauthorized,
		"2FA required",
		"2FA_REQUIRED",
	)
)
