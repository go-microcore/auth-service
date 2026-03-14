// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package tokens

import (
	"go.microcore.dev/framework/transport"
)

var (
	// ErrInvalidAuthRequest indicates the auth request is invalid.
	ErrInvalidAuthRequest = transport.NewError(
		transport.ErrBadRequest,
		"invalid auth request",
		"INVALID_AUTH_REQUEST",
	)

	// ErrRefreshTokenAlreadyUsed indicates that the provided refresh token has
	// already been used.
	ErrRefreshTokenAlreadyUsed = transport.NewError(
		transport.ErrBadRequest,
		"token already used",
		"TOKEN_ALREADY_USED",
	)

	// ErrStaticTokenExist indicates that a static access token with the given
	// ID already exists.
	ErrStaticTokenExist = transport.NewError(
		transport.ErrBadRequest,
		"static token exist",
		"STATIC_TOKEN_EXIST",
	)
)
