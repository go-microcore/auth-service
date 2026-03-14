// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package roles

import (
	"time"
)

type (
	// Data

	// CreateRoleData represents the input data required to create a new role.
	CreateRoleData struct {
		ID          string
		Name        string
		Description string
		SystemFlag  bool
		ServiceFlag bool
	}

	// FilterRolesData represents the criteria for filtering roles.
	FilterRolesData struct {
		ID          *[]string
		Name        *[]string
		SystemFlag  *bool
		ServiceFlag *bool
	}

	// AuthorizeHTTPRolesData represents the data needed to check HTTP role authorization.
	AuthorizeHTTPRolesData struct {
		Roles  []string
		Path   string
		Method string
	}

	// Results

	// CreateRoleResult represents the result returned after creating a role.
	CreateRoleResult struct {
		ID          string
		Name        string
		Description string
		SystemFlag  bool
		ServiceFlag bool
		Created     time.Time
		Updated     time.Time
	}

	// FilterRolesResult represents the result of filtering roles.
	FilterRolesResult struct {
		ID          string
		Name        string
		Description string
		SystemFlag  bool
		ServiceFlag bool
		Created     time.Time
		Updated     time.Time
		HTTPRules   []HTTPRuleResult
	}

	// HTTPRuleResult represents a single HTTP rule associated with a role.
	HTTPRuleResult struct {
		ID      uint
		RoleID  string
		Path    string
		Methods []string
		Mfa     bool
		Created time.Time
		Updated time.Time
	}

	// AuthorizeHTTPRolesResult represents the outcome of checking HTTP role authorization.
	AuthorizeHTTPRolesResult struct {
		Mfa bool
	}
)
