// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package log

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"go.microcore.dev/framework/log"
	"go.microcore.dev/framework/shutdown"
)

// ErrUnknownLogOutput is returned when the log output type is unknown or unsupported.
var ErrUnknownLogOutput = errors.New("unknown log output type")

type (
	// Log wraps the logging configuration and output writer.
	Log struct {
		Config *Config   // Logging configuration
		Stderr io.Writer // Writer for stderr output
	}
)

// Setup initializes the logger based on the configuration.
func (l *Log) Setup() error {
	if err := l.setLogLevel(); err != nil {
		return err
	}

	w, err := l.getWriter()
	if err != nil {
		return err
	}

	return l.configureLogger(w)
}

// setLogLevel sets the log level.
func (l *Log) setLogLevel() error {
	if err := log.SetLevelStr(l.Config.Level); err != nil {
		return newConfigError(fmt.Errorf("set log level: %w", err))
	}

	return nil
}

// getWriter returns the writer based on configuration.
func (l *Log) getWriter() (io.Writer, error) {
	switch l.Config.Output {
	case ConfigOutputStderr:
		return l.Stderr, nil
	case ConfigOutputFile:
		return l.setupFileWriter()
	default:
		return nil, fmt.Errorf("%w: %v", ErrUnknownLogOutput, l.Config.Output)
	}
}

// setupFileWriter opens the log file and registers a shutdown handler.
func (l *Log) setupFileWriter() (io.Writer, error) {
	f, err := os.OpenFile(
		l.Config.File,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		DefaultLogFilePerm,
	)
	if err != nil {
		return nil, newConfigError(fmt.Errorf("open log file %q: %w", l.Config.File, err))
	}

	// Close file on shutdown
	if err := shutdown.AddHandler(func(ctx context.Context, _ int) error {
		errCh := make(chan error, 1)

		go func() { errCh <- f.Close() }()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case e := <-errCh:
			return e
		}
	}); err != nil {
		return nil, fmt.Errorf("add shutdown handler: %w", err)
	}

	return f, nil
}

// configureLogger applies the log configuration.
func (l *Log) configureLogger(w io.Writer) error {
	if err := log.Config(log.Options{
		Writer:      w,
		Format:      l.Config.Format,
		ReplaceAttr: log.DefaultPrettyReplaceAttr,
	}); err != nil {
		return newConfigError(fmt.Errorf("log config: %w", err))
	}

	return nil
}

func newConfigError(err error) error {
	return shutdown.NewExitReason(
		shutdown.ExitConfigError,
		err,
	)
}
