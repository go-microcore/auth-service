// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package log_test

import (
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	clilog "go.microcore.dev/auth-service/cmd/cli/log"
	"go.microcore.dev/framework/log"
)

//nolint:paralleltest // uses global log setup
func TestSetup_SetLevel(t *testing.T) {
	_log := &clilog.Log{
		Config: &clilog.Config{
			Level:  "debug",
			Output: clilog.ConfigOutputStderr,
			Format: log.FormatText,
			File:   "",
		},
		Stderr: os.Stderr,
	}

	err := _log.Setup()
	require.NoError(t, err)
	require.Equal(t, slog.LevelDebug, log.GetLevel())
}

//nolint:paralleltest // uses global log setup
func TestSetup_SetLevelError(t *testing.T) {
	_log := &clilog.Log{
		Config: &clilog.Config{
			Level:  "invalid_level",
			Output: clilog.ConfigOutputStderr,
			Format: log.FormatText,
			File:   "",
		},
		Stderr: os.Stderr,
	}

	err := _log.Setup()
	require.Error(t, err)
	require.ErrorContains(t, err, `set log level: slog: level string "invalid_level": unknown name`)
}

//nolint:paralleltest // uses global log setup
func TestSetup_FormatError(t *testing.T) {
	_log := &clilog.Log{
		Config: &clilog.Config{
			Level:  "debug",
			Output: clilog.ConfigOutputStderr,
			Format: "invalid_format",
			File:   "",
		},
		Stderr: os.Stderr,
	}

	err := _log.Setup()
	require.Error(t, err)
	require.ErrorContains(t, err, `log config: format invalid_format not implemented`)
}

//nolint:paralleltest // uses global log setup
func TestSetup_OutputStderr(t *testing.T) {
	read, write, err := os.Pipe()
	require.NoError(t, err)

	defer read.Close()
	defer write.Close()

	_log := &clilog.Log{
		Config: &clilog.Config{
			Level:  "debug",
			Output: clilog.ConfigOutputStderr,
			Format: log.FormatText,
			File:   "",
		},
		Stderr: write,
	}

	err = _log.Setup()
	require.NoError(t, err)

	log.Info("test message")

	_ = write.Close()

	data, err := io.ReadAll(read)
	require.NoError(t, err)

	require.Contains(t, string(data), `level=INFO msg="test message"`)
}

//nolint:paralleltest // uses global log setup
func TestSetup_OutputFile(t *testing.T) {
	tmpFile := t.TempDir() + "/auth-service.log"

	_log := &clilog.Log{
		Config: &clilog.Config{
			Level:  "debug",
			Output: clilog.ConfigOutputFile,
			Format: log.FormatText,
			File:   tmpFile,
		},
		Stderr: os.Stderr,
	}

	err := _log.Setup()
	require.NoError(t, err)

	log.Info("test message")

	data, err := os.ReadFile(tmpFile)
	require.NoError(t, err)

	require.Contains(t, string(data), `level=INFO msg="test message"`)
}

//nolint:paralleltest // uses global log setup
func TestSetup_OutputFileError(t *testing.T) {
	tmpFile := t.TempDir() + "/nonexistent_dir/auth-service.log"

	_log := &clilog.Log{
		Config: &clilog.Config{
			Level:  "debug",
			Output: clilog.ConfigOutputFile,
			Format: log.FormatText,
			File:   tmpFile,
		},
		Stderr: os.Stderr,
	}

	err := _log.Setup()
	require.Error(t, err)
	require.ErrorContains(t, err, "no such file or directory")
}
