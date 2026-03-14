// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package roles

import (
	"go.microcore.dev/framework/transport"
)

var (
	// ErrRoleNotFound is returned when a role with the specified ID does not exist.
	ErrRoleNotFound = transport.NewError(
		transport.ErrNotFound,
		"role not found",
		"ROLE_NOT_FOUND",
	)

	// ErrRolesInsufficientPermissions is returned when the user does not have
	// sufficient permissions to perform the requested action.
	ErrRolesInsufficientPermissions = transport.NewError(
		transport.ErrForbidden,
		"insufficient permissions",
		"INSUFFICIENT_PERMISSIONS",
	)
)
