// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package rules

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	rulesrp "go.microcore.dev/auth-service/internal/port/adapter/repository/rules"
	rulessp "go.microcore.dev/auth-service/internal/port/service/rules"
)

type (
	// ServiceConfig provides rules service configuration.
	ServiceConfig struct {
		Logger          *slog.Logger
		RulesRepository rulesrp.Adapter
	}

	service struct {
		*ServiceConfig
	}
)

// NewService creates a new instance of the service.
func NewService(config *ServiceConfig) rulessp.Service {
	return &service{config}
}

// CreateHTTPRule creates a new HTTP rule.
func (s *service) CreateHTTPRule(
	ctx context.Context,
	data rulessp.CreateHTTPRuleData,
) (*rulessp.CreateHTTPRuleResult, error) {
	// Check rule exist
	if rules, err := s.RulesRepository.FilterHTTPRules(
		ctx,
		rulesrp.FilterHTTPRulesData{
			ID:      nil,
			RoleID:  &[]string{data.RoleID},
			Path:    &[]string{data.Path},
			Methods: &data.Methods,
			Mfa:     nil,
		},
	); err != nil {
		return nil, fmt.Errorf("filter: %w", err)
	} else if len(rules) == 1 {
		return nil, rulessp.ErrRuleExist
	}

	// Create rule
	rule, err := s.RulesRepository.CreateHTTPRule(
		ctx,
		rulesrp.CreateHTTPRuleData(data),
	)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	// Make res
	res := rulessp.CreateHTTPRuleResult(*rule)

	return &res, nil
}

// FilterHTTPRules returns HTTP rules matching the filter criteria.
func (s *service) FilterHTTPRules(
	ctx context.Context,
	data rulessp.FilterHTTPRulesData,
) ([]rulessp.FilterHTTPRulesResult, error) {
	// Filter rules
	rules, err := s.RulesRepository.FilterHTTPRules(
		ctx,
		rulesrp.FilterHTTPRulesData(data),
	)
	if err != nil {
		return nil, fmt.Errorf("filter: %w", err)
	}

	// Make res
	res := make([]rulessp.FilterHTTPRulesResult, len(rules))
	for i := range rules {
		res[i] = rulessp.FilterHTTPRulesResult{
			ID:      rules[i].ID,
			RoleID:  rules[i].RoleID,
			Path:    rules[i].Path,
			Methods: rules[i].Methods,
			Mfa:     rules[i].Mfa,
			Created: rules[i].Created,
			Updated: rules[i].Updated,
		}
	}

	return res, nil
}

// DeleteHTTPRule deletes an HTTP rule by ID.
func (s *service) DeleteHTTPRule(
	ctx context.Context,
	id uint,
) error {
	// Delete rule
	if err := s.RulesRepository.DeleteHTTPRule(ctx, id); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// UpdateHTTPRule updates an HTTP rule by ID.
func (s *service) UpdateHTTPRule(
	ctx context.Context,
	id uint,
	data map[string]any,
) error {
	// Set rule updated date
	data["updated"] = time.Unix(0, time.Now().UnixNano())

	// Update rule
	if err := s.RulesRepository.UpdateHTTPRule(
		ctx,
		id,
		data,
	); err != nil {
		return fmt.Errorf("update: %w", err)
	}

	return nil
}
