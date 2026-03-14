// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package seed

import (
	"strings"

	"go.microcore.dev/framework/config/env"
)

const (
	// DefaultServiceName is the default name of the seed service.
	DefaultServiceName = "auth-service-seed"

	// DefaultAdminRoleID is the ID for the admin role.
	DefaultAdminRoleID = "admin"

	// DefaultAdminRoleName is the display name for the admin role.
	DefaultAdminRoleName = "Admin"

	// DefaultAdminRoleDescription describes the admin role and its permissions.
	DefaultAdminRoleDescription = "Administrative user role that has full and unrestricted " +
		"permissions to manage all services and system resources. Users with this role have " +
		"complete access to all operations, including critical and system-level functions, " +
		"and can perform any actions without limitations. This role is intended exclusively " +
		"for trusted administrators and must be assigned with extreme caution."

	// DefaultAdminRulePath is the default HTTP rule path for the admin role.
	DefaultAdminRulePath = "/**"

	// DefaultAdminRuleMethods are the default HTTP methods allowed for the admin role.
	DefaultAdminRuleMethods = "GET,HEAD,POST,PUT,PATCH,DELETE"

	// DefaultAdminTokenID is the ID of the admin static token.
	DefaultAdminTokenID = "admin"

	// DefaultAdminTokenDescription describes the admin token and its permissions.
	DefaultAdminTokenDescription = "An administrative system token with full permissions to " +
		"manage all services and resources of the platform. It provides unrestricted access " +
		"to all operations, including critical and system-level actions. The token is issued " +
		"only once, cannot be viewed or recovered again, and must be kept in strict secrecy. " +
		"Compromise of this token results in a complete system compromise."
)

type (
	// Config defines the seed configuration.
	Config struct {
		Name  string
		Admin *ConfigAdmin
	}

	// ConfigAdmin contains admin-related seed configuration.
	ConfigAdmin struct {
		Role  *ConfigAdminRole  // Admin role configuration.
		Rule  *ConfigAdminRule  // Access rule for the admin role.
		Token *ConfigAdminToken // Static admin token configuration.
	}

	// ConfigAdminRole defines the admin role parameters.
	ConfigAdminRole struct {
		ID          string // Unique role identifier.
		Name        string // Human-readable role name.
		Description string // Role description and permissions.
	}

	// ConfigAdminRule defines access rules for the admin role.
	ConfigAdminRule struct {
		Path    string   // Allowed path pattern.
		Methods []string // Allowed HTTP methods.
	}

	// ConfigAdminToken defines the static admin token parameters.
	ConfigAdminToken struct {
		ID          string // Token identifier.
		Description string // Token description and purpose.
	}
)

// NewConfig creates and validates a seed configuration.
func NewConfig() *Config {
	return &Config{
		Name: env.StrDefault("SERVICE_NAME", DefaultServiceName),
		Admin: &ConfigAdmin{
			Role: &ConfigAdminRole{
				ID: env.StrDefault(
					"ADMIN_ROLE_ID",
					DefaultAdminRoleID,
				),
				Name: env.StrDefault(
					"ADMIN_ROLE_NAME",
					DefaultAdminRoleName,
				),
				Description: env.StrDefault(
					"ADMIN_ROLE_DESCRIPTION",
					DefaultAdminRoleDescription,
				),
			},
			Rule: &ConfigAdminRule{
				Path: env.StrDefault(
					"ADMIN_RULE_PATH",
					DefaultAdminRulePath,
				),
				Methods: strings.Split(
					env.StrDefault(
						"ADMIN_RULE_METHODS",
						DefaultAdminRuleMethods,
					), ",",
				),
			},
			Token: &ConfigAdminToken{
				ID: env.StrDefault(
					"ADMIN_TOKEN_ID",
					DefaultAdminTokenID,
				),
				Description: env.StrDefault(
					"ADMIN_TOKEN_DESCRIPTION",
					DefaultAdminTokenDescription,
				),
			},
		},
	}
}
