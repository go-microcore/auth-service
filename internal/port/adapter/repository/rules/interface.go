// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package rules

import (
	"context"
)

type (
	// Adapter defines the interface for the rules repository adapter.
	Adapter interface {
		CreateHTTPRule(
			ctx context.Context,
			data CreateHTTPRuleData,
		) (*CreateHTTPRuleResult, error)

		FilterHTTPRules(
			ctx context.Context,
			data FilterHTTPRulesData,
		) ([]FilterHTTPRulesResult, error)

		DeleteHTTPRule(
			ctx context.Context,
			id uint,
		) error

		UpdateHTTPRule(
			ctx context.Context,
			id uint,
			data map[string]any,
		) error
	}
)
