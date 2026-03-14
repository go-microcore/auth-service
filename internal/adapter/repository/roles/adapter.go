// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package roles

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/gobwas/glob"
	"github.com/lib/pq"
	rolesm "go.microcore.dev/auth-service/internal/adapter/repository/roles/model"
	rulesm "go.microcore.dev/auth-service/internal/adapter/repository/rules/model"
	rolesrp "go.microcore.dev/auth-service/internal/port/adapter/repository/roles"
	"go.microcore.dev/auth-service/internal/shared/postgres"
	"go.microcore.dev/auth-service/internal/shared/redis"
	"gorm.io/gorm"
)

// ErrRolesCacheNotReady signals that the in-memory roles cache has not been initialized yet.
var ErrRolesCacheNotReady = errors.New("roles cache not ready")

type (
	// AdapterConfig provides roles repository adapter configuration.
	AdapterConfig struct {
		Config     *Config
		Logger     *slog.Logger
		Redis      *redis.Redis
		Postgres   *postgres.Postgres
		RolesCache atomic.Value
	}

	adapter struct {
		*AdapterConfig
	}

	// roleRulesCache stores compiled HTTP rules for each role.
	roleRulesCache map[string][]compiledHTTPRule

	// compiledHTTPRule represents a compiled HTTP rule for fast role authorization checks.
	compiledHTTPRule struct {
		Methods map[string]struct{} // HTTP methods allowed
		Path    glob.Glob           // Compiled glob for path matching
		Mfa     bool                // Whether MFA is required
	}
)

// NewAdapter creates a new instance of the adapter.
func NewAdapter(config *AdapterConfig) rolesrp.Adapter {
	return &adapter{config}
}

// CreateRole creates a new role in the database and updates the roles cache.
func (a *adapter) CreateRole(
	ctx context.Context,
	data rolesrp.CreateRoleData,
) (*rolesrp.CreateRoleResult, error) {
	now := time.Now()

	// Create model
	obj := rolesm.AuthRole{
		ID:          data.ID,
		Name:        data.Name,
		Description: data.Description,
		SystemFlag:  data.SystemFlag,
		ServiceFlag: data.ServiceFlag,
		Created:     time.Unix(0, now.UnixNano()),
		Updated:     time.Unix(0, now.UnixNano()),
		HTTPRules:   []rulesm.AuthHTTPRule{},
	}

	// Save role to database
	if err := a.Postgres.Manager.Client().
		WithContext(ctx).
		Create(&obj).
		Error; err != nil {
		return nil, fmt.Errorf("pg: %w", err)
	}

	// Update roles cache
	if err := a.Redis.Manager.Client().
		Publish(ctx, "roles:update", "create_role").
		Err(); err != nil {
		a.Logger.ErrorContext(
			ctx,
			"update roles cache",
			slog.Any("error", err),
			slog.String("action", "create_role"),
		)
	}

	// Make res
	res := rolesrp.CreateRoleResult{
		ID:          obj.ID,
		Name:        obj.Name,
		Description: obj.Description,
		SystemFlag:  obj.SystemFlag,
		ServiceFlag: obj.ServiceFlag,
		Created:     obj.Created,
		Updated:     obj.Updated,
	}

	return &res, nil
}

// FilterRoles retrieves roles from the database according to filter criteria.
func (a *adapter) FilterRoles(
	ctx context.Context,
	data rolesrp.FilterRolesData,
) ([]rolesrp.FilterRolesResult, error) {
	var authRoles []rolesm.AuthRole

	query := a.Postgres.Manager.Client().WithContext(ctx)

	query = a.applyRoleFilters(query, data)

	if err := query.
		Preload("HTTPRules").
		Find(&authRoles).
		Error; err != nil {
		return nil, fmt.Errorf("pg: %w", err)
	}

	return mapAuthRoles(authRoles), nil
}

// DeleteRole deletes a role by ID and updates the roles cache.
func (a *adapter) DeleteRole(
	ctx context.Context,
	id string,
) error {
	// Delete role from database
	result := a.Postgres.Manager.Client().
		WithContext(ctx).
		Delete(&rolesm.AuthRole{}, "id = ?", id)

	// Check errors
	if result.Error != nil {
		return fmt.Errorf("pg: %w", result.Error)
	}

	// If role not found
	if result.RowsAffected == 0 {
		return rolesrp.ErrRoleNotFound
	}

	// Update roles cache
	if err := a.Redis.Manager.Client().
		Publish(ctx, "roles:update", "delete_role").
		Err(); err != nil {
		a.Logger.ErrorContext(
			ctx,
			"update roles cache",
			slog.Any("error", err),
			slog.String("action", "delete_role"),
		)
	}

	return nil
}

// UpdateRole updates role fields by ID and refreshes the roles cache.
func (a *adapter) UpdateRole(
	ctx context.Context,
	id string,
	data map[string]any,
) error {
	// Update role in database
	result := a.Postgres.Manager.Client().
		WithContext(ctx).
		Model(&rolesm.AuthRole{}).Where("id = ?", id).
		Updates(data)

	// Check errors
	if result.Error != nil {
		return fmt.Errorf("pg: %w", result.Error)
	}

	// If role not found
	if result.RowsAffected == 0 {
		return rolesrp.ErrRoleNotFound
	}

	// Update roles cache
	if err := a.Redis.Manager.Client().
		Publish(ctx, "roles:update", "update_role").
		Err(); err != nil {
		a.Logger.ErrorContext(
			ctx,
			"update roles cache",
			slog.Any("error", err),
			slog.String("action", "update_role"),
		)
	}

	return nil
}

// AuthorizeHTTPRoles checks if any of the given roles are authorized to access a specific
// HTTP path and method. It uses the in-memory roles cache for fast lookup and returns the
// MFA requirement if a matching rule is found. If no roles match, it returns
// ErrRolesInsufficientPermissions.
func (a *adapter) AuthorizeHTTPRoles(
	_ context.Context,
	data rolesrp.AuthorizeHTTPRolesData,
) (*rolesrp.AuthorizeHTTPRolesResult, error) {
	cache, ok := a.RolesCache.Load().(roleRulesCache)
	if !ok {
		return nil, ErrRolesCacheNotReady
	}

	for i := range data.Roles {
		rules, ok := cache[data.Roles[i]]
		if !ok {
			continue
		}

		for j := range rules {
			if _, ok := rules[j].Methods[data.Method]; !ok {
				continue
			}

			if rules[j].Path.Match(data.Path) {
				return &rolesrp.AuthorizeHTTPRolesResult{
					Mfa: rules[j].Mfa,
				}, nil
			}
		}
	}

	return nil, rolesrp.ErrRolesInsufficientPermissions
}

// UpdateRolesCache loads all roles and their HTTP rules from the database, compiles them
// into a fast in-memory format, and updates the atomic in-memory cache. This function
// ensures that the roles cache always contains up-to-date rules for fast authorization
// checks.
func (a *adapter) UpdateRolesCache(
	ctx context.Context,
) error {
	// Load roles
	roles, err := a.FilterRoles(
		ctx,
		rolesrp.FilterRolesData{
			ID:          nil,
			Name:        nil,
			SystemFlag:  nil,
			ServiceFlag: nil,
		},
	)
	if err != nil {
		return fmt.Errorf("filter roles: %w", err)
	}

	// Prepare cache data
	compiled := make(roleRulesCache)
	for i := range roles {
		rules := make([]compiledHTTPRule, len(roles[i].HTTPRules))
		for j := range roles[i].HTTPRules {
			methods := make(map[string]struct{}, len(roles[i].HTTPRules[j].Methods))
			for k := range roles[i].HTTPRules[j].Methods {
				methods[roles[i].HTTPRules[j].Methods[k]] = struct{}{}
			}

			rules[j] = compiledHTTPRule{
				Methods: methods,
				Path:    glob.MustCompile(roles[i].HTTPRules[j].Path),
				Mfa:     roles[i].HTTPRules[j].Mfa,
			}
		}

		compiled[roles[i].ID] = rules
	}

	// Set cache data
	a.RolesCache.Store(compiled)

	return nil
}

// SubscribeRoleUpdates subscribes to the "roles:update" Redis channel and listens for role
// update events. When a message is received, it triggers an update of the in-memory roles
// cache. The subscription runs until the provided context is canceled. This ensures that
// all service replicas keep their in-memory cache in sync with role changes.
//
//nolint:gocognit,cyclop,funlen // need refactoring
func (a *adapter) SubscribeRoleUpdates(
	ctx context.Context,
) {
	go func() {
		// Set initial backoff
		backoff := a.Config.Cache.SubRoleUpBackoff.Initial

		for {
			// Shutdown handle
			if ctx.Err() != nil {
				return
			}

			pubsub := a.Redis.Manager.Client().Subscribe(ctx, "roles:update")

			a.Logger.DebugContext(
				ctx,
				"subscribed to redis channel",
				slog.String("channel", "roles:update"),
			)

			for {
				msg, err := pubsub.ReceiveMessage(ctx)
				if err != nil {
					// Shutdown handle
					if ctx.Err() != nil {
						a.Logger.DebugContext(
							ctx,
							"redis subscription stopped by context",
							slog.String("channel", "roles:update"),
						)

						// Close connection
						if err = pubsub.Close(); err != nil {
							a.Logger.ErrorContext(
								ctx,
								"close redis channel",
								slog.String("channel", "roles:update"),
								slog.Any("error", err),
							)
						}

						return
					}

					// Disconnect
					a.Logger.InfoContext(
						ctx,
						"redis subscription reconnecting...",
						slog.String("channel", "roles:update"),
						slog.Any("error", err),
						slog.String("backoff", backoff.String()),
					)

					// Close connection
					if err = pubsub.Close(); err != nil {
						a.Logger.ErrorContext(
							ctx,
							"close redis channel",
							slog.String("channel", "roles:update"),
							slog.Any("error", err),
						)
					}

					break // Reconnect
				}

				// Set initial backoff
				backoff = a.Config.Cache.SubRoleUpBackoff.Initial

				// Message recevied
				a.Logger.DebugContext(
					ctx,
					"received roles update event",
					slog.String("event", msg.Payload),
				)

				// Update roles cache
				if err := a.UpdateRolesCache(ctx); err != nil {
					a.Logger.ErrorContext(
						ctx,
						"update roles cache",
						slog.Any("error", err),
					)
				} else {
					a.Logger.DebugContext(
						ctx,
						"roles cache updated successfully",
					)
				}
			}

			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return
			}

			if backoff < a.Config.Cache.SubRoleUpBackoff.Max {
				backoff = time.Duration(
					int64(backoff) * a.Config.Cache.SubRoleUpBackoff.Multiplier,
				)
			}
		}
	}()
}

// PeriodicRolesCacheSync periodically refreshes the in-memory roles cache by reloading all
// roles and their HTTP rules from the database. This ensures that even if a PUB/SUB event
// was missed while the service was offline or the Redis connection was down, the cache
// remains up-to-date. The function runs until the provided context is canceled, and the
// refresh interval is controlled by the `interval` parameter.
func (a *adapter) PeriodicRolesCacheSync(
	ctx context.Context,
) {
	go func() {
		ticker := time.NewTicker(a.Config.Cache.PeriodicRolesSyncInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := a.UpdateRolesCache(ctx); err != nil {
					a.Logger.ErrorContext(
						ctx,
						"periodic roles cache sync",
						slog.Any("error", err),
					)
				} else {
					a.Logger.DebugContext(
						ctx,
						"periodic roles cache sync successful",
					)
				}

			case <-ctx.Done():
				return
			}
		}
	}()
}

func (*adapter) applyRoleFilters(
	query *gorm.DB,
	data rolesrp.FilterRolesData,
) *gorm.DB {
	if data.ID != nil {
		query = query.Where("id IN ?", *data.ID)
	}

	if data.Name != nil {
		query = query.Where("name IN ?", *data.Name)
	}

	if data.SystemFlag != nil {
		query = query.Where("system_flag = ?", *data.SystemFlag)
	}

	if data.ServiceFlag != nil {
		query = query.Where("service_flag = ?", *data.ServiceFlag)
	}

	return query
}

func mapAuthRoles(
	roles []rolesm.AuthRole,
) []rolesrp.FilterRolesResult {
	res := make([]rolesrp.FilterRolesResult, len(roles))
	for i := range roles {
		httpRules := make([]rolesrp.HTTPRuleResult, len(roles[i].HTTPRules))
		for j := range roles[i].HTTPRules {
			httpRules[j] = rolesrp.HTTPRuleResult{
				ID:      roles[i].HTTPRules[j].ID,
				RoleID:  roles[i].HTTPRules[j].RoleID,
				Path:    roles[i].HTTPRules[j].Path,
				Methods: pq.StringArray(roles[i].HTTPRules[j].Methods), //nolint:unconvert // need convert
				Mfa:     roles[i].HTTPRules[j].Mfa,
				Created: roles[i].HTTPRules[j].Created,
				Updated: roles[i].HTTPRules[j].Updated,
			}
		}

		res[i] = rolesrp.FilterRolesResult{
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

	return res
}
