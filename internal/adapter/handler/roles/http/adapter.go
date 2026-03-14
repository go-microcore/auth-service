// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package http

import (
	"context"
	"fmt"
	"log/slog"

	roleshp "go.microcore.dev/auth-service/internal/port/adapter/handler/roles/http"
	rolessp "go.microcore.dev/auth-service/internal/port/service/roles"
	"go.microcore.dev/framework/transport/http"
	"go.microcore.dev/framework/transport/http/server"
)

type (
	// AdapterConfig provides roles HTTP adapter handler configuration.
	AdapterConfig struct {
		Logger       *slog.Logger
		RolesService rolessp.Service
	}

	adapter struct {
		*AdapterConfig
	}
)

// NewAdapter creates a new instance of the adapter.
func NewAdapter(config *AdapterConfig) roleshp.Adapter {
	return &adapter{config}
}

// CreateRole creates a new role.
//
// @Summary Creates a new role.
// @Tags Roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body roleshp.CreateRoleRequest true "Request data"
// @Success 201 {object} roleshp.CreateRoleResponse
// @Failure 400 {object} server.ErrResponse "Possible codes: INVALID_ROLE_ID, INVALID_ROLE_NAME, INVALID_ROLE_DESCRIPTION, ROLE_EXIST_ID"
// @Failure 401 {object} server.ErrResponse "Possible codes: INVALID_TOKEN, 2FA_REQUIRED"
// @Failure 403 {object} server.ErrResponse "Possible codes: INSUFFICIENT_PERMISSIONS"
// @Failure 415 {object} server.ErrResponse "Possible codes: INVALID_JSON_BODY"
// @Failure 503 {object} server.ErrResponse "Possible codes: SERVICE_UNAVAILABLE"
// @Router /auth/roles/ [post]
func (a *adapter) CreateRole(
	ctx context.Context,
	c *server.RequestContext,
	req *roleshp.CreateRoleRequest,
) {
	// Create role
	res, err := a.RolesService.CreateRole(
		ctx,
		rolessp.CreateRoleData(*req),
	)
	if err != nil {
		c.WriteError(fmt.Errorf("create role: %w", err))
		return
	}

	// Write res
	c.WriteJsonWithStatusCode(
		http.StatusCreated,
		roleshp.CreateRoleResponse(*res),
	)
}

// FilterRoles returns roles matching the filter criteria.
//
// @Summary Returns roles matching the filter criteria.
// @Tags Roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body roleshp.FilterRolesRequest true "Request data"
// @Success 200 {array} roleshp.FilterRolesResponse
// @Failure 401 {object} server.ErrResponse "Possible codes: INVALID_TOKEN, 2FA_REQUIRED"
// @Failure 403 {object} server.ErrResponse "Possible codes: INSUFFICIENT_PERMISSIONS"
// @Failure 415 {object} server.ErrResponse "Possible codes: INVALID_JSON_BODY"
// @Failure 503 {object} server.ErrResponse "Possible codes: SERVICE_UNAVAILABLE"
// @Router /auth/roles/filter [post]
func (a *adapter) FilterRoles(
	ctx context.Context,
	c *server.RequestContext,
	req *roleshp.FilterRolesRequest,
) {
	// Filter roles
	roles, err := a.RolesService.FilterRoles(
		ctx,
		rolessp.FilterRolesData(*req),
	)
	if err != nil {
		c.WriteError(fmt.Errorf("filter roles: %w", err))
		return
	}

	// Make res
	res := make([]roleshp.FilterRolesResponse, len(roles))
	for i := range roles {
		httpRules := make([]roleshp.HTTPRuleResponse, len(roles[i].HTTPRules))
		for j := range roles[i].HTTPRules {
			httpRules[j] = roleshp.HTTPRuleResponse(roles[i].HTTPRules[j])
		}

		res[i] = roleshp.FilterRolesResponse{
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

	// Write res
	c.WriteJsonWithStatusCode(http.StatusOK, res)
}

// UpdateRole updates role data by ID.
//
// @Summary Updates role data by ID.
// @Tags Roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Param request body roleshp._UpdateRoleRequest true "Request data"
// @Success 204
// @Failure 400 {object} server.ErrResponse "Possible codes: INVALID_ROLE_ID, INVALID_ROLE_NAME, INVALID_ROLE_DESCRIPTION"
// @Failure 401 {object} server.ErrResponse "Possible codes: INVALID_TOKEN, 2FA_REQUIRED"
// @Failure 403 {object} server.ErrResponse "Possible codes: INSUFFICIENT_PERMISSIONS"
// @Failure 404 {object} server.ErrResponse "Possible codes: ROLE_NOT_FOUND"
// @Failure 415 {object} server.ErrResponse "Possible codes: INVALID_JSON_BODY"
// @Failure 503 {object} server.ErrResponse "Possible codes: SERVICE_UNAVAILABLE"
// @Router /auth/roles/{id} [patch]
func (a *adapter) UpdateRole(
	ctx context.Context,
	c *server.RequestContext,
	req *roleshp.UpdateRoleRequest,
) {
	// Get role id
	id, err := c.UserValueStr("id")
	if err != nil {
		c.WriteError(roleshp.ErrUserRoleInvalidID)
		return
	}

	role := make(map[string]any)

	if req.ID.Set {
		role["id"] = req.ID.Value
	}

	if req.Name.Set {
		role["name"] = req.Name.Value
	}

	if req.Description.Set {
		role["description"] = req.Description.Value
	}

	if req.SystemFlag.Set {
		role["system_flag"] = req.SystemFlag.Value
	}

	if req.ServiceFlag.Set {
		role["service_flag"] = req.ServiceFlag.Value
	}

	// Update role
	if err := a.RolesService.UpdateRole(
		ctx,
		id,
		role,
	); err != nil {
		c.WriteError(fmt.Errorf("update role: %w", err))
		return
	}

	// Write res
	c.StatusCode(http.StatusNoContent)
}

// DeleteRole deletes a role by ID.
//
// @Summary Deletes a role by ID.
// @Tags Roles
// @Security BearerAuth
// @Produce json
// @Param id path string true "Role ID"
// @Success 204
// @Failure 400 {object} server.ErrResponse "Possible codes: INVALID_ROLE_ID"
// @Failure 401 {object} server.ErrResponse "Possible codes: INVALID_TOKEN, 2FA_REQUIRED"
// @Failure 403 {object} server.ErrResponse "Possible codes: INSUFFICIENT_PERMISSIONS"
// @Failure 404 {object} server.ErrResponse "Possible codes: ROLE_NOT_FOUND"
// @Failure 503 {object} server.ErrResponse "Possible codes: SERVICE_UNAVAILABLE"
// @Router /auth/roles/{id} [delete]
func (a *adapter) DeleteRole(
	ctx context.Context,
	c *server.RequestContext,
) {
	// Get role id
	id, err := c.UserValueStr("id")
	if err != nil {
		c.WriteError(roleshp.ErrUserRoleInvalidID)
		return
	}

	// Delete role
	if err := a.RolesService.DeleteRole(ctx, id); err != nil {
		c.WriteError(fmt.Errorf("delete role: %w", err))
		return
	}

	// Write res
	c.StatusCode(http.StatusNoContent)
}
