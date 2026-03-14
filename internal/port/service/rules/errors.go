// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package rules

import (
	"go.microcore.dev/framework/transport"
)

// ErrRuleExist indicates that an HTTP rule with the same parameters
// already exists.
var ErrRuleExist = transport.NewError(
	transport.ErrBadRequest,
	"rule exist",
	"RULE_EXIST",
)
