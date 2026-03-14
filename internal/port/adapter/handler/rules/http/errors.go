// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package http

import (
	"go.microcore.dev/framework/transport"
)

var (
	// ErrHTTPRuleInvalidID indicates the ID is invalid.
	ErrHTTPRuleInvalidID = transport.NewError(
		transport.ErrBadRequest,
		"invalid id",
		"INVALID_ID",
	)

	// ErrHTTPRuleInvalidRoleID indicates the RoleID is invalid.
	ErrHTTPRuleInvalidRoleID = transport.NewError(
		transport.ErrBadRequest,
		"invalid role id",
		"INVALID_ROLE_ID",
	)

	// ErrHTTPRuleInvalidPath indicates the Path is invalid.
	ErrHTTPRuleInvalidPath = transport.NewError(
		transport.ErrBadRequest,
		"invalid path",
		"INVALID_PATH",
	)

	// ErrHTTPRuleInvalidMethods indicates the Methods slice is invalid
	// or empty.
	ErrHTTPRuleInvalidMethods = transport.NewError(
		transport.ErrBadRequest,
		"invalid methods",
		"INVALID_METHODS",
	)

	// ErrHTTPRuleInvalidMfa indicates the Mfa value is invalid.
	ErrHTTPRuleInvalidMfa = transport.NewError(
		transport.ErrBadRequest,
		"invalid mfa",
		"INVALID_MFA",
	)
)
