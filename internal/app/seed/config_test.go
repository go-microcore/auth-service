// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package seed_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.microcore.dev/auth-service/internal/app/seed"
)

func TestNewConfig(t *testing.T) {
	t.Setenv("SERVICE_NAME", "service")
	t.Setenv("ADMIN_ROLE_ID", "id")
	t.Setenv("ADMIN_ROLE_NAME", "name")
	t.Setenv("ADMIN_ROLE_DESCRIPTION", "descr")
	t.Setenv("ADMIN_RULE_PATH", "path")
	t.Setenv("ADMIN_RULE_METHODS", "get,post")
	t.Setenv("ADMIN_TOKEN_ID", "id")
	t.Setenv("ADMIN_TOKEN_DESCRIPTION", "descr")

	cfg := seed.NewConfig()

	require.Equal(t, "service", cfg.Name)
	require.Equal(t, "id", cfg.Admin.Role.ID)
	require.Equal(t, "name", cfg.Admin.Role.Name)
	require.Equal(t, "descr", cfg.Admin.Role.Description)
	require.Equal(t, "path", cfg.Admin.Rule.Path)
	require.Equal(t, []string{"get", "post"}, cfg.Admin.Rule.Methods)
	require.Equal(t, "id", cfg.Admin.Token.ID)
	require.Equal(t, "descr", cfg.Admin.Token.Description)
}

//nolint:paralleltest // uses environment variables
func TestNewConfig_Default(t *testing.T) {
	cfg := seed.NewConfig()

	require.Equal(t, seed.DefaultServiceName, cfg.Name)
	require.Equal(t, seed.DefaultAdminRoleID, cfg.Admin.Role.ID)
	require.Equal(t, seed.DefaultAdminRoleName, cfg.Admin.Role.Name)
	require.Equal(t, seed.DefaultAdminRoleDescription, cfg.Admin.Role.Description)
	require.Equal(t, seed.DefaultAdminRulePath, cfg.Admin.Rule.Path)
	require.Equal(t, strings.Split(seed.DefaultAdminRuleMethods, ","), cfg.Admin.Rule.Methods)
	require.Equal(t, seed.DefaultAdminTokenID, cfg.Admin.Token.ID)
	require.Equal(t, seed.DefaultAdminTokenDescription, cfg.Admin.Token.Description)
}
