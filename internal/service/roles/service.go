// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package roles

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	rolesrp "go.microcore.dev/auth-service/internal/port/adapter/repository/roles"
	rolessp "go.microcore.dev/auth-service/internal/port/service/roles"
)

type (
	// ServiceConfig provides roles service configuration.
	ServiceConfig struct {
		Logger          *slog.Logger
		RolesRepository rolesrp.Adapter
	}

	service struct {
		*ServiceConfig
	}
)

// NewService creates a new instance of the service.
func NewService(config *ServiceConfig) rolessp.Service {
	return &service{config}
}

// CreateRole creates a new role.
func (s *service) CreateRole(
	ctx context.Context,
	data rolessp.CreateRoleData,
) (*rolessp.CreateRoleResult, error) {
	// Check role exist
	if roles, err := s.RolesRepository.FilterRoles(
		ctx,
		rolesrp.FilterRolesData{
			ID:          &[]string{data.ID},
			Name:        nil,
			SystemFlag:  nil,
			ServiceFlag: nil,
		},
	); err != nil {
		return nil, fmt.Errorf("filter: %w", err)
	} else if len(roles) == 1 {
		return nil, rolessp.ErrRoleExistID
	}

	// Create role
	role, err := s.RolesRepository.CreateRole(
		ctx,
		rolesrp.CreateRoleData(data),
	)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	// Make res
	res := rolessp.CreateRoleResult(*role)

	return &res, nil
}

// FilterRoles returns roles matching the filter criteria.
func (s *service) FilterRoles(
	ctx context.Context,
	data rolessp.FilterRolesData,
) ([]rolessp.FilterRolesResult, error) {
	// Filter roles
	roles, err := s.RolesRepository.FilterRoles(
		ctx,
		rolesrp.FilterRolesData(data),
	)
	if err != nil {
		return nil, fmt.Errorf("filter: %w", err)
	}

	// Make res
	res := make([]rolessp.FilterRolesResult, len(roles))
	for i := range roles {
		httpRules := make([]rolessp.HTTPRuleResult, len(roles[i].HTTPRules))
		for j := range roles[i].HTTPRules {
			httpRules[j] = rolessp.HTTPRuleResult(roles[i].HTTPRules[j])
		}

		res[i] = rolessp.FilterRolesResult{
			ID:          roles[i].ID,
			Name:        roles[i].Name,
			Description: roles[i].Description,
			SystemFlag:  roles[i].SystemFlag,
			ServiceFlag: roles[i].ServiceFlag,
			Created:     roles[i].Created,
			Updated:     roles[i].Updated,
			HTTPRules:   httpRules,
		}
	}

	return res, nil
}

// DeleteRole deletes a role by ID.
func (s *service) DeleteRole(
	ctx context.Context,
	id string,
) error {
	// Delete role
	if err := s.RolesRepository.DeleteRole(ctx, id); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// UpdateRole updates role data by ID.
func (s *service) UpdateRole(
	ctx context.Context,
	id string,
	data map[string]any,
) error {
	// Set role updated date
	data["updated"] = time.Unix(0, time.Now().UnixNano())

	// Update role
	if err := s.RolesRepository.UpdateRole(
		ctx,
		id,
		data,
	); err != nil {
		return fmt.Errorf("update: %w", err)
	}

	return nil
}

// AuthorizeHTTPRoles checks HTTP authorization for given roles.
func (s *service) AuthorizeHTTPRoles(
	ctx context.Context,
	data rolessp.AuthorizeHTTPRolesData,
) (*rolessp.AuthorizeHTTPRolesResult, error) {
	// Authorize HTTP roles
	result, err := s.RolesRepository.AuthorizeHTTPRoles(
		ctx,
		rolesrp.AuthorizeHTTPRolesData(data),
	)
	if err != nil {
		return nil, fmt.Errorf("authorize: %w", err)
	}

	// Make res
	res := rolessp.AuthorizeHTTPRolesResult(*result)

	return &res, nil
}
