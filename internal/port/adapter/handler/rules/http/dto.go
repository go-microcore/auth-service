// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package http

import (
	"time"

	"go.microcore.dev/sdk/types"
)

type (
	// Requests

	// CreateHTTPRuleRequest represents the payload to create a new HTTP rule.
	CreateHTTPRuleRequest struct {
		RoleID  string   `json:"roleId"`
		Path    string   `json:"path"`
		Methods []string `json:"methods"`
		Mfa     bool     `json:"mfa"`
	}

	// FilterHTTPRulesRequest represents filters for querying HTTP rules.
	FilterHTTPRulesRequest struct {
		ID      *[]uint   `json:"id"`
		RoleID  *[]string `json:"roleId"`
		Path    *[]string `json:"path"`
		Methods *[]string `json:"methods"`
		Mfa     *bool     `json:"mfa"`
	}

	// _UpdateHTTPRuleRequest is used internally for Swagger documentation.
	//
	//nolint:unused // for Swagger
	_UpdateHTTPRuleRequest struct {
		RoleID  uint     `json:"roleId"`
		Path    string   `json:"path"`
		Methods []string `json:"methods"`
		Mfa     bool     `json:"mfa"`
	}

	// UpdateHTTPRuleRequest represents the payload to update an HTTP rule.
	UpdateHTTPRuleRequest struct {
		RoleID  types.Nullable[string]   `json:"roleId"`
		Path    types.Nullable[string]   `json:"path"`
		Methods types.Nullable[[]string] `json:"methods"`
		Mfa     types.Nullable[bool]     `json:"mfa"`
	}

	// Responses

	// CreateHTTPRuleResponse represents the HTTP rule returned after creation.
	CreateHTTPRuleResponse struct {
		ID      uint      `json:"id"`
		RoleID  string    `json:"roleId"`
		Path    string    `json:"path"`
		Methods []string  `json:"methods"`
		Mfa     bool      `json:"mfa"`
		Created time.Time `json:"created"`
		Updated time.Time `json:"updated"`
	}

	// FilterHTTPRulesResponse represents a single HTTP rule returned by a filter query.
	FilterHTTPRulesResponse struct {
		ID      uint      `json:"id"`
		RoleID  string    `json:"roleId"`
		Path    string    `json:"path"`
		Methods []string  `json:"methods"`
		Mfa     bool      `json:"mfa"`
		Created time.Time `json:"created"`
		Updated time.Time `json:"updated"`
	}
)

// Validate validates the CreateHTTPRuleRequest fields.
func (r *CreateHTTPRuleRequest) Validate() error {
	if err := r.ValidateRoleID(); err != nil {
		return err
	}

	if err := r.ValidatePath(); err != nil {
		return err
	}

	return r.ValidateMethods()
}

// ValidateRoleID checks that RoleID is not empty.
func (r *CreateHTTPRuleRequest) ValidateRoleID() error {
	if r.RoleID == "" {
		return ErrHTTPRuleInvalidRoleID
	}

	return nil
}

// ValidatePath checks that Path is not empty.
func (r *CreateHTTPRuleRequest) ValidatePath() error {
	if r.Path == "" {
		return ErrHTTPRuleInvalidPath
	}

	return nil
}

// ValidateMethods checks that Methods is not empty.
func (r *CreateHTTPRuleRequest) ValidateMethods() error {
	if len(r.Methods) == 0 {
		return ErrHTTPRuleInvalidMethods
	}

	return nil
}

// Validate checks all fields of UpdateHTTPRuleRequest.
func (r *UpdateHTTPRuleRequest) Validate() error {
	if err := r.ValidateRoleID(); err != nil {
		return err
	}

	if err := r.ValidatePath(); err != nil {
		return err
	}

	if err := r.ValidateMethods(); err != nil {
		return err
	}

	return r.ValidateMfa()
}

// ValidateRoleID checks if RoleID is set and non-empty.
func (r *UpdateHTTPRuleRequest) ValidateRoleID() error {
	if r.RoleID.Set && (r.RoleID.Value == nil || *r.RoleID.Value == "") {
		return ErrHTTPRuleInvalidRoleID
	}

	return nil
}

// ValidatePath checks if Path is set and non-empty.
func (r *UpdateHTTPRuleRequest) ValidatePath() error {
	if r.Path.Set && (r.Path.Value == nil || *r.Path.Value == "") {
		return ErrHTTPRuleInvalidPath
	}

	return nil
}

// ValidateMethods checks if Methods is set and contains at least one element.
func (r *UpdateHTTPRuleRequest) ValidateMethods() error {
	if r.Methods.Set && (r.Methods.Value == nil || len(*r.Methods.Value) == 0) {
		return ErrHTTPRuleInvalidMethods
	}

	return nil
}

// ValidateMfa checks if Mfa is set and not nil.
func (r *UpdateHTTPRuleRequest) ValidateMfa() error {
	if r.Mfa.Set && r.Mfa.Value == nil {
		return ErrHTTPRuleInvalidMfa
	}

	return nil
}
