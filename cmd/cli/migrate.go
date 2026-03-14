// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package cli

import (
	"github.com/spf13/cobra"
	"go.microcore.dev/auth-service/internal/app/migrate"
)

func newMigrateCmd() *cobra.Command {
	return newCmd(
		"migrate",
		"Run migrate",
		migrate.Init,
		&migrate.Options{},
		nil,
	)
}
