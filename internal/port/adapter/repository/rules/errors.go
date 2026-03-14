// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package rules

import (
	"go.microcore.dev/framework/transport"
)

// ErrRuleNotFound indicates that the requested HTTP rule does not exist.
var ErrRuleNotFound = transport.NewError(
	transport.ErrNotFound,
	"rule not found",
	"RULE_NOT_FOUND",
)
