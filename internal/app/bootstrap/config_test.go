// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package bootstrap_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.microcore.dev/auth-service/internal/app/bootstrap"
)

func Test_NewConfig(t *testing.T) {
	t.Setenv("SERVICE_NAME", "service")

	cfg := bootstrap.NewConfig()

	require.Equal(t, "service", cfg.Name)
}

//nolint:paralleltest // uses environment variables
func Test_NewConfig_Default(t *testing.T) {
	cfg := bootstrap.NewConfig()

	require.Equal(t, bootstrap.DefaultServiceName, cfg.Name)
}
