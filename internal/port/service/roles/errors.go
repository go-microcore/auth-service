// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package roles

import (
	"go.microcore.dev/framework/transport"
)

var (
	// ErrRoleNotFound is returned when a role with the specified ID or name
	// does not exist.
	ErrRoleNotFound = transport.NewError(
		transport.ErrBadRequest,
		"role not found",
		"ROLE_NOT_FOUND",
	)

	// ErrRoleExistID is returned when trying to create a role with an ID that
	// already exists.
	ErrRoleExistID = transport.NewError(
		transport.ErrBadRequest,
		"role exist id",
		"ROLE_EXIST_ID",
	)

	// ErrRoleExistName is returned when trying to create a role with a name that
	// already exists.
	ErrRoleExistName = transport.NewError(
		transport.ErrBadRequest,
		"role exist name",
		"ROLE_EXIST_NAME",
	)
)
