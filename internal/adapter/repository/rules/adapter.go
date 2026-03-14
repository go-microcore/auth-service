// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package rules

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/lib/pq"
	"go.microcore.dev/auth-service/internal/adapter/repository/rules/model"
	rulesrp "go.microcore.dev/auth-service/internal/port/adapter/repository/rules"
	"go.microcore.dev/auth-service/internal/shared/postgres"
	"go.microcore.dev/auth-service/internal/shared/redis"
)

type (
	// AdapterConfig provides rules repository adapter configuration.
	AdapterConfig struct {
		Logger   *slog.Logger
		Redis    *redis.Redis
		Postgres *postgres.Postgres
	}

	adapter struct {
		*AdapterConfig
	}
)

// NewAdapter creates a new instance of the adapter.
func NewAdapter(config *AdapterConfig) rulesrp.Adapter {
	return &adapter{config}
}

// CreateHTTPRule creates a new HTTP rule for a role and updates the roles cache.
func (a *adapter) CreateHTTPRule(
	ctx context.Context,
	data rulesrp.CreateHTTPRuleData,
) (*rulesrp.CreateHTTPRuleResult, error) {
	now := time.Now()

	// Create model
	obj := model.AuthHTTPRule{
		RoleID:  data.RoleID,
		Path:    data.Path,
		Methods: pq.StringArray(data.Methods),
		Mfa:     data.Mfa,
		Created: time.Unix(0, now.UnixNano()),
		Updated: time.Unix(0, now.UnixNano()),
	}

	// Save rule to database
	if err := a.Postgres.Manager.Client().
		WithContext(ctx).
		Create(&obj).
		Error; err != nil {
		return nil, fmt.Errorf("pg: %w", err)
	}

	// Update roles cache
	if err := a.Redis.Manager.Client().
		Publish(ctx, "roles:update", "create_rule").
		Err(); err != nil {
		a.Logger.ErrorContext(
			ctx,
			"update roles cache",
			slog.Any("error", err),
			slog.String("action", "create_rule"),
		)
	}

	// Make res
	res := rulesrp.CreateHTTPRuleResult{
		ID:      obj.ID,
		RoleID:  obj.RoleID,
		Path:    obj.Path,
		Methods: []string(obj.Methods),
		Mfa:     obj.Mfa,
		Created: obj.Created,
		Updated: obj.Updated,
	}

	return &res, nil
}

// FilterHTTPRules retrieves HTTP rules from the database based on filter criteria.
func (a *adapter) FilterHTTPRules(
	ctx context.Context,
	data rulesrp.FilterHTTPRulesData,
) ([]rulesrp.FilterHTTPRulesResult, error) {
	// Create model
	rules := []model.AuthHTTPRule{}

	// Create query with context
	query := a.Postgres.Manager.Client().WithContext(ctx)

	// Filter by id
	if data.ID != nil {
		query = query.Where("id IN ?", *data.ID)
	}

	// Filter by role_id
	if data.RoleID != nil {
		query = query.Where("role_id IN ?", *data.RoleID)
	}

	// Filter by path
	if data.Path != nil {
		query = query.Where("path IN ?", *data.Path)
	}

	// Filter by methods (exact match)
	if data.Methods != nil {
		query = query.Where(
			"methods @> ? AND methods <@ ?",
			pq.Array(*data.Methods),
			pq.Array(*data.Methods),
		)
	}

	// Filter by mfa
	if data.Mfa != nil {
		query = query.Where("mfa = ?", *data.Mfa)
	}

	// Get rules from database
	if err := query.Find(&rules).Error; err != nil {
		return nil, fmt.Errorf("pg: %w", err)
	}

	// Make res
	res := make([]rulesrp.FilterHTTPRulesResult, len(rules))
	for i := range rules {
		res[i] = rulesrp.FilterHTTPRulesResult{
			ID:      rules[i].ID,
			RoleID:  rules[i].RoleID,
			Path:    rules[i].Path,
			Methods: []string(rules[i].Methods),
			Mfa:     rules[i].Mfa,
			Created: rules[i].Created,
			Updated: rules[i].Updated,
		}
	}

	return res, nil
}

// DeleteHTTPRule deletes an HTTP rule by ID and updates the roles cache.
func (a *adapter) DeleteHTTPRule(
	ctx context.Context,
	id uint,
) error {
	// Delete rule from database
	result := a.Postgres.Manager.Client().
		WithContext(ctx).
		Delete(&model.AuthHTTPRule{}, "id = ?", id)

	// Check errors
	if result.Error != nil {
		return fmt.Errorf("pg: %w", result.Error)
	}

	// If rule not found
	if result.RowsAffected == 0 {
		return rulesrp.ErrRuleNotFound
	}

	// Update roles cache
	if err := a.Redis.Manager.Client().
		Publish(ctx, "roles:update", "delete_rule").
		Err(); err != nil {
		a.Logger.ErrorContext(
			ctx,
			"update roles cache",
			slog.Any("error", err),
			slog.String("action", "delete_rule"),
		)
	}

	return nil
}

// UpdateHTTPRule updates an existing HTTP rule by ID and refreshes the roles cache.
func (a *adapter) UpdateHTTPRule(
	ctx context.Context,
	id uint,
	data map[string]any,
) error {
	// Update rule in database
	result := a.Postgres.Manager.Client().
		WithContext(ctx).
		Model(&model.AuthHTTPRule{}).
		Where("id = ?", id).Updates(data)

	// Check errors
	if result.Error != nil {
		return fmt.Errorf("pg: %w", result.Error)
	}

	// If rule not found
	if result.RowsAffected == 0 {
		return rulesrp.ErrRuleNotFound
	}

	// Update roles cache
	if err := a.Redis.Manager.Client().
		Publish(ctx, "roles:update", "update_rule").
		Err(); err != nil {
		a.Logger.ErrorContext(
			ctx,
			"update roles cache",
			slog.Any("error", err),
			slog.String("action", "update_rule"),
		)
	}

	return nil
}
