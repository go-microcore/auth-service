// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package log_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	clilog "go.microcore.dev/auth-service/cmd/cli/log"
	"go.microcore.dev/framework/log"
)

func TestNewConfig_OutputStderr(t *testing.T) {
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("LOG_OUTPUT", string(clilog.ConfigOutputStderr))
	t.Setenv("LOG_FORMAT", string(log.FormatText))

	cfg, err := clilog.NewConfig()
	require.NoError(t, err)

	require.Equal(t, "info", cfg.Level)
	require.Equal(t, clilog.ConfigOutputStderr, cfg.Output)
	require.Equal(t, log.FormatText, cfg.Format)
}

func TestNewConfig_OutputFile(t *testing.T) {
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("LOG_OUTPUT", string(clilog.ConfigOutputFile))
	t.Setenv("LOG_FORMAT", string(log.FormatText))
	t.Setenv("LOG_FILE", "/var/log/auth-service.log")

	cfg, err := clilog.NewConfig()
	require.NoError(t, err)

	require.Equal(t, "info", cfg.Level)
	require.Equal(t, clilog.ConfigOutputFile, cfg.Output)
	require.Equal(t, log.FormatText, cfg.Format)
	require.Equal(t, "/var/log/auth-service.log", cfg.File)
}

func TestNewConfig_OutputFileError(t *testing.T) {
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("LOG_OUTPUT", string(clilog.ConfigOutputFile))
	t.Setenv("LOG_FORMAT", string(log.FormatText))

	_, err := clilog.NewConfig()
	require.Error(t, err)
	require.ErrorContains(t, err, "variable LOG_FILE is not set")
}

func TestNewConfig_OutputError(t *testing.T) {
	t.Setenv("LOG_OUTPUT", "invalid_output")

	_, err := clilog.NewConfig()
	require.Error(t, err)
	require.ErrorContains(t, err, `log output not implemented: "invalid_output"`)
}

//nolint:paralleltest // uses environment variables
func TestNewConfig_Default(t *testing.T) {
	cfg, err := clilog.NewConfig()
	require.NoError(t, err)

	require.Equal(t, clilog.DefaultLogLevel, cfg.Level)
	require.Equal(t, clilog.DefaultLogOutput, string(cfg.Output))
	require.Equal(t, clilog.DefaultLogFormat, string(cfg.Format))
}
