// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package cli

import (
	"os"

	"github.com/spf13/cobra"
	"go.microcore.dev/auth-service/internal/app/bootstrap"
)

func newBootstrapCmd() *cobra.Command {
	return newCmd(
		"bootstrap",
		"Run bootstrap",
		bootstrap.Init,
		&bootstrap.Options{
			Output: &bootstrap.OutputOptions{
				JSON:          false,
				JSONPretty:    false,
				Quiet:         false,
				JwtAccessKey:  "",
				JwtRefreshKey: "",
				JwtHashKey:    "",
				AuthKey:       "",
				Stdout:        os.Stdout,
			},
		},
		bootstrapFlagConfig,
	)
}

func bootstrapFlagConfig(cmd *cobra.Command, options *bootstrap.Options) {
	cmd.Flags().BoolVarP(
		&options.Output.JSON,
		"json",
		"j",
		false,
		"output in JSON format to stdout, suppressing all non-JSON output "+
			"(overrides --quiet, --jwt-access-key-output, --jwt-refresh-key-output, "+
			"--jwt-hash-key-output, --auth-key-output)",
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
		"suppress security warning and output only the keys value",
	)

	cmd.Flags().StringVar(
		&options.Output.JwtAccessKey,
		"jwt-access-key-output",
		"stdout",
		"where to output the jwt access key: use 'stdout' or 'file:/path/to/file.txt'")

	cmd.Flags().StringVar(
		&options.Output.JwtRefreshKey,
		"jwt-refresh-key-output",
		"stdout",
		"where to output the jwt refresh key: use 'stdout' or 'file:/path/to/file.txt'")

	cmd.Flags().StringVar(
		&options.Output.JwtHashKey,
		"jwt-hash-key-output",
		"stdout",
		"where to output the jwt hash key: use 'stdout' or 'file:/path/to/file.txt'")

	cmd.Flags().StringVar(
		&options.Output.AuthKey,
		"auth-key-output",
		"stdout",
		"where to output the auth key: use 'stdout' or 'file:/path/to/file.txt'")

	cmd.MarkFlagsMutuallyExclusive("json", "quiet")
	cmd.MarkFlagsMutuallyExclusive("json", "jwt-access-key-output")
	cmd.MarkFlagsMutuallyExclusive("json", "jwt-refresh-key-output")
	cmd.MarkFlagsMutuallyExclusive("json", "jwt-hash-key-output")
	cmd.MarkFlagsMutuallyExclusive("json", "auth-key-output")
}
