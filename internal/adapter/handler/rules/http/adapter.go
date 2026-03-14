// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package http

import (
	"context"
	"fmt"
	"log/slog"

	ruleshp "go.microcore.dev/auth-service/internal/port/adapter/handler/rules/http"
	rulessp "go.microcore.dev/auth-service/internal/port/service/rules"
	"go.microcore.dev/framework/transport/http"
	"go.microcore.dev/framework/transport/http/server"
)

type (
	// AdapterConfig provides rules HTTP adapter handler configuration.
	AdapterConfig struct {
		Logger       *slog.Logger
		RulesService rulessp.Service
	}

	adapter struct {
		*AdapterConfig
	}
)

// NewAdapter creates a new instance of the adapter.
func NewAdapter(config *AdapterConfig) ruleshp.Adapter {
	return &adapter{config}
}

// CreateHTTPRule creates a new HTTP rule.
//
// @Summary Creates a new HTTP rule.
// @Tags Rules (HTTP)
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body ruleshp.CreateHTTPRuleRequest true "Request data"
// @Success 201 {object} ruleshp.CreateHTTPRuleResponse
// @Failure 400 {object} server.ErrResponse "Possible codes: INVALID_ROLE_ID, INVALID_PATH, INVALID_METHODS, RULE_EXIST"
// @Failure 401 {object} server.ErrResponse "Possible codes: INVALID_TOKEN, 2FA_REQUIRED"
// @Failure 403 {object} server.ErrResponse "Possible codes: INSUFFICIENT_PERMISSIONS"
// @Failure 415 {object} server.ErrResponse "Possible codes: INVALID_JSON_BODY"
// @Failure 503 {object} server.ErrResponse "Possible codes: SERVICE_UNAVAILABLE"
// @Router /auth/rules/http/ [post]
func (a *adapter) CreateHTTPRule(
	ctx context.Context,
	c *server.RequestContext,
	req *ruleshp.CreateHTTPRuleRequest,
) {
	// Create HTTP rule
	res, err := a.RulesService.CreateHTTPRule(
		ctx,
		rulessp.CreateHTTPRuleData(*req),
	)
	if err != nil {
		c.WriteError(fmt.Errorf("create http rule: %w", err))
		return
	}

	// Write res
	c.WriteJsonWithStatusCode(
		http.StatusCreated,
		ruleshp.CreateHTTPRuleResponse(*res),
	)
}

// FilterHTTPRules returns HTTP rules matching the filter criteria.
//
// @Summary Returns HTTP rules matching the filter criteria.
// @Tags Rules (HTTP)
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body ruleshp.FilterHTTPRulesRequest true "Request data"
// @Success 200 {array} ruleshp.FilterHTTPRulesResponse
// @Failure 401 {object} server.ErrResponse "Possible codes: INVALID_TOKEN, 2FA_REQUIRED"
// @Failure 403 {object} server.ErrResponse "Possible codes: INSUFFICIENT_PERMISSIONS"
// @Failure 415 {object} server.ErrResponse "Possible codes: INVALID_JSON_BODY"
// @Failure 503 {object} server.ErrResponse "Possible codes: SERVICE_UNAVAILABLE"
// @Router /auth/rules/http/filter [post]
func (a *adapter) FilterHTTPRules(
	ctx context.Context,
	c *server.RequestContext,
	req *ruleshp.FilterHTTPRulesRequest,
) {
	rules, err := a.RulesService.FilterHTTPRules(ctx, rulessp.FilterHTTPRulesData(*req))
	if err != nil {
		c.WriteError(fmt.Errorf("filter http rules: %w", err))
		return
	}

	// Make res
	res := make([]ruleshp.FilterHTTPRulesResponse, len(rules))
	for i := range rules {
		res[i] = ruleshp.FilterHTTPRulesResponse(rules[i])
	}

	// Write res
	c.WriteJsonWithStatusCode(http.StatusOK, res)
}

// UpdateHTTPRule updates an HTTP rule by ID.
//
// @Summary Updates an HTTP rule by ID.
// @Tags Rules (HTTP)
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Rule ID"
// @Param request body ruleshp._UpdateHTTPRuleRequest true "Request data"
// @Success 204
// @Failure 400 {object} server.ErrResponse "Possible codes: INVALID_ID, INVALID_ROLE_ID, INVALID_PATH, INVALID_METHODS, INVALID_MFA"
// @Failure 401 {object} server.ErrResponse "Possible codes: INVALID_TOKEN, 2FA_REQUIRED"
// @Failure 403 {object} server.ErrResponse "Possible codes: INSUFFICIENT_PERMISSIONS"
// @Failure 404 {object} server.ErrResponse "Possible codes: RULE_NOT_FOUND"
// @Failure 415 {object} server.ErrResponse "Possible codes: INVALID_JSON_BODY"
// @Failure 503 {object} server.ErrResponse "Possible codes: SERVICE_UNAVAILABLE"
// @Router /auth/rules/http/{id} [patch]
func (a *adapter) UpdateHTTPRule(
	ctx context.Context,
	c *server.RequestContext,
	req *ruleshp.UpdateHTTPRuleRequest,
) {
	// Get rule id
	id, err := c.UserValueUint("id")
	if err != nil {
		c.WriteError(ruleshp.ErrHTTPRuleInvalidID)
		return
	}

	rule := make(map[string]any)

	if req.RoleID.Set {
		rule["role_id"] = req.RoleID.Value
	}

	if req.Path.Set {
		rule["path"] = req.Path.Value
	}

	if req.Methods.Set {
		rule["methods"] = req.Methods.Value
	}

	if req.Mfa.Set {
		rule["mfa"] = req.Mfa.Value
	}

	// Update HTTP rule
	if err := a.RulesService.UpdateHTTPRule(
		ctx,
		id,
		rule,
	); err != nil {
		c.WriteError(fmt.Errorf("update http rule: %w", err))
		return
	}

	c.StatusCode(http.StatusNoContent)
}

// DeleteHTTPRule deletes an HTTP rule by ID.
//
// @Summary Deletes an HTTP rule by ID.
// @Tags Rules (HTTP)
// @Security BearerAuth
// @Produce json
// @Param id path int true "Rule ID"
// @Success 204
// @Failure 400 {object} server.ErrResponse "Possible codes: INVALID_ID"
// @Failure 401 {object} server.ErrResponse "Possible codes: INVALID_TOKEN, 2FA_REQUIRED"
// @Failure 403 {object} server.ErrResponse "Possible codes: INSUFFICIENT_PERMISSIONS"
// @Failure 404 {object} server.ErrResponse "Possible codes: RULE_NOT_FOUND"
// @Failure 503 {object} server.ErrResponse "Possible codes: SERVICE_UNAVAILABLE"
// @Router /auth/rules/http/{id} [delete]
func (a *adapter) DeleteHTTPRule(
	ctx context.Context,
	c *server.RequestContext,
) {
	// Get rule id
	id, err := c.UserValueUint("id")
	if err != nil {
		c.WriteError(ruleshp.ErrHTTPRuleInvalidID)
		return
	}

	// Delete HTTP rule
	if err := a.RulesService.DeleteHTTPRule(ctx, id); err != nil {
		c.WriteError(fmt.Errorf("delete http rule: %w", err))
		return
	}

	// Write res
	c.StatusCode(http.StatusNoContent)
}
