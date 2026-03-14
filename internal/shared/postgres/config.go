// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package postgres

import (
	"go.microcore.dev/framework/config/env"
)

const (
	// DefaultPostgresHost is the default Postgres host.
	DefaultPostgresHost = "postgres"

	// DefaultPostgresPort is the default Postgres port.
	DefaultPostgresPort = "5432"

	// DefaultPostgresUser is the default Postgres user.
	DefaultPostgresUser = "postgres"

	// DefaultPostgresPassword is the default Postgres password.
	DefaultPostgresPassword = "password"

	// DefaultPostgresDB is the default Postgres database name.
	DefaultPostgresDB = "postgres"

	// DefaultPostgresSSL is the default Postgres SSL mode.
	DefaultPostgresSSL = "disable"
)

type (
	// Config defines the Postgres configuration.
	Config struct {
		Host     string
		Port     string
		User     string
		Password string
		DB       string
		SSL      string
	}
)

// NewConfig creates and validates a Postgres configuration.
func NewConfig() *Config {
	return &Config{
		Host:     env.StrDefault("POSTGRES_HOST", DefaultPostgresHost),
		Port:     env.StrDefault("POSTGRES_PORT", DefaultPostgresPort),
		User:     env.StrDefault("POSTGRES_USER", DefaultPostgresUser),
		Password: env.StrDefault("POSTGRES_PASSWORD", DefaultPostgresPassword),
		DB:       env.StrDefault("POSTGRES_DB", DefaultPostgresDB),
		SSL:      env.StrDefault("POSTGRES_SSL", DefaultPostgresSSL),
	}
}
