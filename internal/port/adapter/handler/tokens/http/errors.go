// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package http

import (
	"go.microcore.dev/framework/transport"
)

var (
	// ErrInvalidToken indicates the token provided is invalid.
	ErrInvalidToken = transport.NewError(
		transport.ErrUnauthorized,
		"invalid token",
		"INVALID_TOKEN",
	)

	// ErrInvalidPath indicates the path provided is invalid.
	ErrInvalidPath = transport.NewError(
		transport.ErrBadRequest,
		"invalid path",
		"INVALID_PATH",
	)

	// ErrInvalidMethod indicates the HTTP method provided is invalid.
	ErrInvalidMethod = transport.NewError(
		transport.ErrBadRequest,
		"invalid method",
		"INVALID_METHOD",
	)

	// ErrInvalidRoles indicates that no roles or invalid roles were provided.
	ErrInvalidRoles = transport.NewError(
		transport.ErrBadRequest,
		"invalid roles",
		"INVALID_ROLES",
	)

	// ErrInvalidID indicates the ID provided is invalid.
	ErrInvalidID = transport.NewError(
		transport.ErrBadRequest,
		"invalid id",
		"INVALID_ID",
	)

	// ErrInvalidDescription indicates the description provided is invalid.
	ErrInvalidDescription = transport.NewError(
		transport.ErrBadRequest,
		"invalid description",
		"INVALID_DESCRIPTION",
	)
)
