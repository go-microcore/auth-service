// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package cli

import (
	"os"

	"github.com/spf13/cobra"
	"go.microcore.dev/auth-service/internal/app/seed"
)

func newSeedCmd() *cobra.Command {
	return newCmd(
		"seed",
		"Run seed",
		seed.Init,
		&seed.Options{
			Output: &seed.OutputOptions{
				JSON:       false,
				JSONPretty: false,
				Quiet:      false,
				AdminToken: "",
				Stdout:     os.Stdout,
			},
		},
		seedFlagConfig,
	)
}

func seedFlagConfig(cmd *cobra.Command, options *seed.Options) {
	cmd.Flags().BoolVarP(
		&options.Output.JSON,
		"json",
		"j",
		false,
		"output in JSON format to stdout, suppressing all non-JSON output "+
			"(overrides --quiet and --admin-token-output)",
	)

	cmd.Flags().BoolVarP(
		&options.Output.JSONPretty,
		"json-pretty",
		"p",
		false,
		"output in JSON format to stdout with pretty printing (only applies if --json is set)",
	)

	cmd.Flags().BoolVarP(
		&options.Output.Quiet,
		"quiet",
		"q",
		false,
		"suppress security warning and output only the token value",
	)

	cmd.Flags().StringVar(
		&options.Output.AdminToken,
		"admin-token-output",
		"stdout",
		"where to output the admin token: use 'stdout' or 'file:/path/to/file.txt'")

	cmd.MarkFlagsMutuallyExclusive("json", "quiet")
	cmd.MarkFlagsMutuallyExclusive("json", "admin-token-output")
}
