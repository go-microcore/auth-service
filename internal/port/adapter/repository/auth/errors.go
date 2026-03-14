// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package auth

import (
	"go.microcore.dev/framework/transport"
)

var (
	// ErrInvalidToken is returned when a JWT or access token is invalid.
	ErrInvalidToken = transport.NewError(
		transport.ErrUnauthorized,
		"invalid token",
		"INVALID_TOKEN",
	)

	// ErrStaticTokenNotFound is returned when a requested static access
	// token does not exist.
	ErrStaticTokenNotFound = transport.NewError(
		transport.ErrNotFound,
		"static token not found",
		"STATIC_TOKEN_NOT_FOUND",
	)

	// ErrSessionNotFound is returned when the requested session does not exist.
	ErrSessionNotFound = transport.NewError(
		transport.ErrNotFound,
		"session not found",
		"SESSION_NOT_FOUND",
	)
)
