// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package roles

import (
	"time"
)

type (
	// Data

	// CreateRoleData holds the input data for creating a new role.
	CreateRoleData struct {
		ID          string
		Name        string
		Description string
		SystemFlag  bool
		ServiceFlag bool
	}

	// FilterRolesData holds optional filters for querying roles.
	FilterRolesData struct {
		ID          *[]string
		Name        *[]string
		SystemFlag  *bool
		ServiceFlag *bool
	}

	// AuthorizeHTTPRolesData contains information for checking HTTP role access.
	AuthorizeHTTPRolesData struct {
		Roles  []string
		Path   string
		Method string
	}

	// Results

	// CreateRoleResult represents the result of a role creation operation.
	CreateRoleResult struct {
		ID          string
		Name        string
		Description string
		SystemFlag  bool
		ServiceFlag bool
		Created     time.Time
		Updated     time.Time
	}

	// FilterRolesResult represents a role returned from a filter query.
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

	// HTTPRuleResult represents an HTTP rule associated with a role.
	HTTPRuleResult struct {
		ID      uint
		RoleID  string
		Path    string
		Methods []string
		Mfa     bool
		Created time.Time
		Updated time.Time
	}

	// AuthorizeHTTPRolesResult indicates whether a request passed HTTP role authorization.
	AuthorizeHTTPRolesResult struct {
		Mfa bool
	}
)
