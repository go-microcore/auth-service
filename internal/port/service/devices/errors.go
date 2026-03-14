// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package devices

import (
	"go.microcore.dev/framework/transport"
)

// ErrRefreshTokenAlreadyUsed indicates that the refresh token has already been
// used and cannot be reused.
var ErrRefreshTokenAlreadyUsed = transport.NewError(
	transport.ErrBadRequest,
	"token already used",
	"TOKEN_ALREADY_USED",
)
