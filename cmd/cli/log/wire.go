//go:build wireinject

// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package log

import (
	"io"

	"github.com/google/wire"
)

type Options struct {
	Stderr io.Writer
}

func Init(opts *Options) (*Log, error) {
	wire.Build(
		NewConfig,
		newLog,
	)
	return nil, nil
}

func newLog(
	opts *Options,
	config *Config,
) *Log {
	return &Log{
		Config: config,
		Stderr: opts.Stderr,
	}
}
