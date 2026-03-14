// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package migrate_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.microcore.dev/auth-service/internal/app/migrate"
)

func TestNewConfig(t *testing.T) {
	t.Setenv("SERVICE_NAME", "service")

	cfg := migrate.NewConfig()

	require.Equal(t, "service", cfg.Name)
}

//nolint:paralleltest // uses environment variables
func TestNewConfig_Default(t *testing.T) {
	cfg := migrate.NewConfig()

	require.Equal(t, migrate.DefaultServiceName, cfg.Name)
}
