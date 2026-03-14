// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package log

import (
	"errors"
	"fmt"
	"slices"

	"go.microcore.dev/framework/config/env"
	"go.microcore.dev/framework/log"
)

const (
	// Default

	// DefaultLogLevel is the default severity level for logs.
	DefaultLogLevel = "info"

	// DefaultLogOutput is the default output destination for logs.
	DefaultLogOutput = "stderr"

	// DefaultLogFormat is the default format for log messages.
	DefaultLogFormat = "pretty"

	// DefaultLogFilePerm is the default file permission for log files.
	DefaultLogFilePerm = 0o644

	// Config

	// ConfigOutputStderr writes logs to the standard error stream.
	ConfigOutputStderr ConfigOutput = "stderr"

	// ConfigOutputFile writes logs to a file specified in Config.File.
	ConfigOutputFile ConfigOutput = "file"
)

// ErrLogOutputNotImplemented occurs when the configured log output type is not supported.
var ErrLogOutputNotImplemented = errors.New("log output not implemented")

type (
	// Config defines the logs configuration.
	Config struct {
		Level  string
		Output ConfigOutput
		Format log.OutputFormat
		File   string
	}

	// ConfigOutput represents the output destination for logs.
	ConfigOutput string
)

// NewConfig creates and validates a logs configuration.
func NewConfig() (*Config, error) {
	config := &Config{
		Level:  env.StrDefault("LOG_LEVEL", DefaultLogLevel),
		Output: ConfigOutput(env.StrDefault("LOG_OUTPUT", DefaultLogOutput)),
		Format: log.OutputFormat(env.StrDefault("LOG_FORMAT", DefaultLogFormat)),
		File:   "",
	}

	if !slices.Contains([]ConfigOutput{
		ConfigOutputStderr,
		ConfigOutputFile,
	}, config.Output) {
		return nil, newConfigError(
			fmt.Errorf("%w: %q", ErrLogOutputNotImplemented, config.Output),
		)
	}

	if config.Output == ConfigOutputFile {
		val, err := env.Str("LOG_FILE")
		if err != nil {
			return nil, newConfigError(err)
		}

		config.File = val
	}

	return config, nil
}
