// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package rules

import (
	"time"
)

type (
	// Data

	// CreateHTTPRuleData represents the data required to create a new HTTP rule.
	CreateHTTPRuleData struct {
		RoleID  string
		Path    string
		Methods []string
		Mfa     bool
	}

	// FilterHTTPRulesData represents the filtering criteria for querying HTTP rules.
	FilterHTTPRulesData struct {
		ID      *[]uint
		RoleID  *[]string
		Path    *[]string
		Methods *[]string
		Mfa     *bool
	}

	// Results

	// CreateHTTPRuleResult represents the result returned after creating an HTTP rule.
	CreateHTTPRuleResult struct {
		ID      uint
		RoleID  string
		Path    string
		Methods []string
		Mfa     bool
		Created time.Time
		Updated time.Time
	}

	// FilterHTTPRulesResult represents an HTTP rule returned from a filter/query operation.
	FilterHTTPRulesResult struct {
		ID      uint
		RoleID  string
		Path    string
		Methods []string
		Mfa     bool
		Created time.Time
		Updated time.Time
	}
)
