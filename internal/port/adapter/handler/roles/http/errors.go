// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package http

import (
	"go.microcore.dev/framework/transport"
)

var (
	// ErrUserRoleInvalidID indicates that the role ID is invalid or missing.
	ErrUserRoleInvalidID = transport.NewError(
		transport.ErrBadRequest,
		"invalid role id",
		"INVALID_ROLE_ID",
	)

	// ErrUserRoleInvalidName indicates that the role name is invalid or missing.
	ErrUserRoleInvalidName = transport.NewError(
		transport.ErrBadRequest,
		"invalid role name",
		"INVALID_ROLE_NAME",
	)

	// ErrUserRoleInvalidDescription indicates that the role description is
	// invalid or missing.
	ErrUserRoleInvalidDescription = transport.NewError(
		transport.ErrBadRequest,
		"invalid role description",
		"INVALID_ROLE_DESCRIPTION",
	)
)
