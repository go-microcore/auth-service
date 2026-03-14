// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package cli

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	clilog "go.microcore.dev/auth-service/cmd/cli/log"
	"go.microcore.dev/framework/config/env"
	"go.microcore.dev/framework/log"
	"go.microcore.dev/framework/shutdown"
)

type (
	// Runner represents a CLI or application component that can be executed
	// with a context, returning an error if something goes wrong.
	Runner interface {
		Run(ctx context.Context) error
	}
)

//nolint:gochecknoglobals // set with -ldflags on build.
var (
	// Version is the current version of the application.
	Version = "dev"

	// Commit is the git commit hash at the time of build.
	Commit = "none"

	// BuildDate is the timestamp (RFC3339) when the binary was built.
	BuildDate = "unknown"
)

// Run executes the CLI application, handling shutdown and exit codes.
// It returns a numeric exit code suitable for process termination.
func Run() int {
	defer shutdown.Recover()

	// Create root context
	ctx, err := shutdown.NewContext()
	if err != nil {
		return shutdown.ExitGeneralError
	}

	// Execute root cmd
	if reason := NewRootCmd().ExecuteContext(ctx); reason != nil {
		code, err := shutdown.ParseExitReason(reason)
		if !err {
			return code
		}

		log.ErrorContext(
			ctx,
			reason.Error(),
			slog.Int("code", code),
		)

		return code
	}

	return shutdown.ExitOK
}

// NewRootCmd creates the root CLI command with persistent flags and subcommands.
func NewRootCmd() *cobra.Command {
	var envPath string

	cmd := &cobra.Command{
		Use:                "auth-service",
		Short:              "Auth service for microservices with JWT, RBAC and session management",
		Version:            Version,
		SilenceErrors:      true,
		SilenceUsage:       true,
		DisableSuggestions: true,
		DisableAutoGenTag:  true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if err := initEnv(cmd, envPath); err != nil {
				return shutdown.NewExitReason(
					shutdown.ExitConfigError,
					err,
				)
			}

			if err := initLogs(); err != nil {
				return shutdown.NewExitReason(
					shutdown.ExitConfigError,
					err,
				)
			}

			log.InfoContext(
				cmd.Context(),
				"app",
				slog.String("version", Version),
				slog.String("commit", Commit),
				slog.String("date", BuildDate),
			)

			return nil
		},
	}

	cmd.SetVersionTemplate(fmt.Sprintf(
		"Version: %s\nCommit: %s\nBuildDate: %s\n\n",
		Version, Commit, BuildDate,
	))

	cmd.PersistentFlags().StringVarP(
		&envPath,
		"env",
		"e",
		"",
		"path to .env file",
	)

	cmd.AddCommand(
		newBootstrapCmd(),
		newMigrateCmd(),
		newSeedCmd(),
		newServerCmd(),
	)

	return cmd
}

func initEnv(cmd *cobra.Command, envPath string) error {
	if !cmd.Flags().Changed("env") {
		return nil
	}

	if err := env.New(envPath); err != nil {
		return fmt.Errorf("load env file: %w", err)
	}

	return nil
}

func initLogs() error {
	logs, err := clilog.Init(
		&clilog.Options{
			Stderr: os.Stderr,
		},
	)
	if err != nil {
		return fmt.Errorf("init log: %w", err)
	}

	if err := logs.Setup(); err != nil {
		return fmt.Errorf("setup log: %w", err)
	}

	return nil
}

func newExitSoftwareReason(e error) error {
	return shutdown.NewExitReason(shutdown.ExitSoftware, e)
}

func newCmd[T any, R Runner](
	use string,
	short string,
	initFn func(ctx context.Context, opts *T) (R, error),
	opts *T,
	flagConfig func(cmd *cobra.Command, opts *T),
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, _ []string) error {
			app, err := initFn(cmd.Context(), opts)
			if err != nil {
				return newExitSoftwareReason(fmt.Errorf("init app: %w", err))
			}

			if err := app.Run(cmd.Context()); err != nil {
				return newExitSoftwareReason(fmt.Errorf("run app: %w", err))
			}

			return nil
		},
	}

	if flagConfig != nil {
		flagConfig(cmd, opts)
	}

	return cmd
}
