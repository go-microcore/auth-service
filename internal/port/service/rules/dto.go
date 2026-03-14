// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package rules

import (
	"time"
)

type (
	// Data

	// CreateHTTPRuleData holds the input data needed to create a new HTTP rule.
	CreateHTTPRuleData struct {
		RoleID  string   // Role ID associated with the rule
		Path    string   // HTTP path for the rule
		Methods []string // Allowed HTTP methods
		Mfa     bool     // Whether MFA is required
	}

	// FilterHTTPRulesData holds optional filters for querying HTTP rules.
	FilterHTTPRulesData struct {
		ID      *[]uint   // Optional list of rule IDs
		RoleID  *[]string // Optional list of role IDs
		Path    *[]string // Optional list of paths
		Methods *[]string // Optional list of HTTP methods
		Mfa     *bool     // Optional MFA requirement filter
	}

	// Results

	// CreateHTTPRuleResult represents the result returned after creating an HTTP rule.
	CreateHTTPRuleResult struct {
		ID      uint      // Rule ID
		RoleID  string    // Associated role ID
		Path    string    // HTTP path
		Methods []string  // Allowed HTTP methods
		Mfa     bool      // MFA requirement
		Created time.Time // Creation timestamp
		Updated time.Time // Last updated timestamp
	}

	// FilterHTTPRulesResult represents the result of a filtered HTTP rules query.
	FilterHTTPRulesResult struct {
		ID      uint      // Rule ID
		RoleID  string    // Associated role ID
		Path    string    // HTTP path
		Methods []string  // Allowed HTTP methods
		Mfa     bool      // MFA requirement
		Created time.Time // Creation timestamp
		Updated time.Time // Last updated timestamp
	}
)
