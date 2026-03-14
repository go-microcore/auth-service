//go:build wireinject

// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package seed

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/wire"
	authra "go.microcore.dev/auth-service/internal/adapter/repository/auth"
	rolesra "go.microcore.dev/auth-service/internal/adapter/repository/roles"
	rulesra "go.microcore.dev/auth-service/internal/adapter/repository/rules"
	authrp "go.microcore.dev/auth-service/internal/port/adapter/repository/auth"
	rolesrp "go.microcore.dev/auth-service/internal/port/adapter/repository/roles"
	rulesrp "go.microcore.dev/auth-service/internal/port/adapter/repository/rules"
	rolessp "go.microcore.dev/auth-service/internal/port/service/roles"
	rulessp "go.microcore.dev/auth-service/internal/port/service/rules"
	tokenssp "go.microcore.dev/auth-service/internal/port/service/tokens"
	rolessa "go.microcore.dev/auth-service/internal/service/roles"
	rulessa "go.microcore.dev/auth-service/internal/service/rules"
	tokenssa "go.microcore.dev/auth-service/internal/service/tokens"
	sharedhttp "go.microcore.dev/auth-service/internal/shared/http/client"
	sharedpg "go.microcore.dev/auth-service/internal/shared/postgres"
	sharedredis "go.microcore.dev/auth-service/internal/shared/redis"
	sharedtel "go.microcore.dev/auth-service/internal/shared/telemetry"
	"go.microcore.dev/framework/log"
)

func Init(ctx context.Context, opts *Options) (*Seed, error) {
	wire.Build(
		NewConfig,
		newLogger,
		newTelemetry,
		newPostgres,
		newRedis,
		newHTTPClient,
		// Repository
		newAuthRepository,
		newRolesRepository,
		newRulesRepository,
		// Services
		newTokensService,
		newRolesService,
		newRulesService,

		newSeed,
	)
	return nil, nil
}

func newLogger(
	cfg *Config,
) *slog.Logger {
	return log.New(cfg.Name)
}

func newTelemetry(
	ctx context.Context,
	opts *Options,
	cfg *Config,
) (*sharedtel.Telemetry, error) {
	t, err := sharedtel.Init(
		ctx,
		&sharedtel.Options{
			AppName: cfg.Name,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("init telemetry: %w", err)
	}
	return t, nil
}

func newPostgres(
	telemetry *sharedtel.Telemetry,
) (*sharedpg.Postgres, error) {
	p, err := sharedpg.Init(
		&sharedpg.Options{
			Telemetry: telemetry,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("init postgres: %w", err)
	}
	if err := p.Setup(); err != nil {
		return nil, fmt.Errorf("setup postgres: %w", err)
	}
	return p, nil
}

func newRedis(
	telemetry *sharedtel.Telemetry,
) (*sharedredis.Redis, error) {
	p, err := sharedredis.Init(
		&sharedredis.Options{
			Telemetry: telemetry,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("init redis: %w", err)
	}
	if err := p.Setup(); err != nil {
		return nil, fmt.Errorf("setup redis: %w", err)
	}
	return p, nil
}

func newHTTPClient(
	telemetry *sharedtel.Telemetry,
) (*sharedhttp.Client, error) {
	c, err := sharedhttp.Init(
		&sharedhttp.Options{
			Telemetry: telemetry,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("init http client: %w", err)
	}
	c.Setup()
	return c, nil
}

func newAuthRepository(
	postgres *sharedpg.Postgres,
	redis *sharedredis.Redis,
) (authrp.Adapter, error) {
	r, err := authra.Init(
		&authra.Options{
			Postgres: postgres,
			Redis:    redis,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("init auth repository: %w", err)
	}
	return r, nil
}

func newRolesRepository(
	postgres *sharedpg.Postgres,
	redis *sharedredis.Redis,
) (rolesrp.Adapter, error) {
	r, err := rolesra.Init(
		&rolesra.Options{
			Postgres: postgres,
			Redis:    redis,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("init roles repository: %w", err)
	}
	return r, nil
}

func newRulesRepository(
	postgres *sharedpg.Postgres,
	redis *sharedredis.Redis,
) (rulesrp.Adapter, error) {
	r, err := rulesra.Init(
		&rulesra.Options{
			Postgres: postgres,
			Redis:    redis,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("init rules repository: %w", err)
	}
	return r, nil
}

func newTokensService(
	authRepository authrp.Adapter,
) (tokenssp.Service, error) {
	s, err := tokenssa.Init(
		&tokenssa.Options{
			AuthRepository: authRepository,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("init tokens service: %w", err)
	}
	return s, nil
}

func newRulesService(
	rulesRepository rulesrp.Adapter,
) (rulessp.Service, error) {
	s, err := rulessa.Init(
		&rulessa.Options{
			RulesRepository: rulesRepository,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("init rules service: %w", err)
	}
	return s, nil
}

func newRolesService(
	rolesRepository rolesrp.Adapter,
) (rolessp.Service, error) {
	s, err := rolessa.Init(
		&rolessa.Options{
			RolesRepository: rolesRepository,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("init roles service: %w", err)
	}
	return s, nil
}

func newSeed(
	options *Options,
	config *Config,
	logger *slog.Logger,
	telemetry *sharedtel.Telemetry,
	postgres *sharedpg.Postgres,
	redis *sharedredis.Redis,
	httpClient *sharedhttp.Client,
	tokensService tokenssp.Service,
	rolesService rolessp.Service,
	rulesService rulessp.Service,
) *Seed {
	return &Seed{
		Options:    options,
		Config:     config,
		Logger:     logger,
		Telemetry:  telemetry,
		Postgres:   postgres,
		Redis:      redis,
		HTTPClient: httpClient,
		Service: &Service{
			Tokens: tokensService,
			Roles:  rolesService,
			Rules:  rulesService,
		},
	}
}
