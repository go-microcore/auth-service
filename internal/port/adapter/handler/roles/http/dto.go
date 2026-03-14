// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package http

import (
	"time"

	"go.microcore.dev/sdk/types"
)

type (
	// Requests

	// CreateRoleRequest represents the request payload for creating a new role.
	CreateRoleRequest struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		SystemFlag  bool   `json:"systemFlag"`
		ServiceFlag bool   `json:"serviceFlag"`
	}

	// FilterRolesRequest represents optional filters for querying roles.
	FilterRolesRequest struct {
		ID          *[]string `json:"id"`
		Name        *[]string `json:"name"`
		SystemFlag  *bool     `json:"systemFlag"`
		ServiceFlag *bool     `json:"serviceFlag"`
	}

	// _UpdateRoleRequest is an internal struct used for documentation (Swagger).
	//
	//nolint:unused // for Swagger
	_UpdateRoleRequest struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	// UpdateRoleRequest represents the request payload for updating an existing role.
	UpdateRoleRequest struct {
		ID          types.Nullable[string] `json:"id"`
		Name        types.Nullable[string] `json:"name"`
		Description types.Nullable[string] `json:"description"`
		SystemFlag  types.Nullable[bool]   `json:"systemFlag"`
		ServiceFlag types.Nullable[bool]   `json:"serviceFlag"`
	}

	// Responses

	// CreateRoleResponse represents the response payload after creating a role.
	CreateRoleResponse struct {
		ID          string    `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		SystemFlag  bool      `json:"systemFlag"`
		ServiceFlag bool      `json:"serviceFlag"`
		Created     time.Time `json:"created"`
		Updated     time.Time `json:"updated"`
	}

	// FilterRolesResponse represents the response payload for a role query.
	FilterRolesResponse struct {
		ID          string             `json:"id"`
		Name        string             `json:"name"`
		Description string             `json:"description"`
		SystemFlag  bool               `json:"systemFlag"`
		ServiceFlag bool               `json:"serviceFlag"`
		Created     time.Time          `json:"created"`
		Updated     time.Time          `json:"updated"`
		HTTPRules   []HTTPRuleResponse `json:"httpRules"`
	}

	// HTTPRuleResponse represents a single HTTP rule associated with a role.
	HTTPRuleResponse struct {
		ID      uint      `json:"id"`
		RoleID  string    `json:"roleId"`
		Path    string    `json:"path"`
		Methods []string  `json:"methods"`
		Mfa     bool      `json:"mfa"`
		Created time.Time `json:"created"`
		Updated time.Time `json:"updated"`
	}
)

// Validate checks that all required fields of CreateRoleRequest are set.
func (r *CreateRoleRequest) Validate() error {
	if r.ID == "" {
		return ErrUserRoleInvalidID
	}

	if r.Name == "" {
		return ErrUserRoleInvalidName
	}

	if r.Description == "" {
		return ErrUserRoleInvalidDescription
	}

	return nil
}

// Validate checks that all set fields in UpdateRoleRequest are valid.
func (r *UpdateRoleRequest) Validate() error {
	if r.ID.Set && (r.ID.Value == nil || *r.ID.Value == "") {
		return ErrUserRoleInvalidID
	}

	if r.Name.Set && (r.Name.Value == nil || *r.Name.Value == "") {
		return ErrUserRoleInvalidName
	}

	if r.Description.Set && (r.Description.Value == nil || *r.Description.Value == "") {
		return ErrUserRoleInvalidDescription
	}

	return nil
}
