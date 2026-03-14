//go:build wireinject

// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/wire"
	authha "go.microcore.dev/auth-service/internal/adapter/handler/auth/http"
	devicesha "go.microcore.dev/auth-service/internal/adapter/handler/devices/http"
	logoutha "go.microcore.dev/auth-service/internal/adapter/handler/logout/http"
	rolesha "go.microcore.dev/auth-service/internal/adapter/handler/roles/http"
	rulesha "go.microcore.dev/auth-service/internal/adapter/handler/rules/http"
	tokensha "go.microcore.dev/auth-service/internal/adapter/handler/tokens/http"
	authra "go.microcore.dev/auth-service/internal/adapter/repository/auth"
	rolesra "go.microcore.dev/auth-service/internal/adapter/repository/roles"
	rulesra "go.microcore.dev/auth-service/internal/adapter/repository/rules"
	authhp "go.microcore.dev/auth-service/internal/port/adapter/handler/auth/http"
	deviceshp "go.microcore.dev/auth-service/internal/port/adapter/handler/devices/http"
	logouthp "go.microcore.dev/auth-service/internal/port/adapter/handler/logout/http"
	roleshp "go.microcore.dev/auth-service/internal/port/adapter/handler/roles/http"
	ruleshp "go.microcore.dev/auth-service/internal/port/adapter/handler/rules/http"
	tokenshp "go.microcore.dev/auth-service/internal/port/adapter/handler/tokens/http"
	authrp "go.microcore.dev/auth-service/internal/port/adapter/repository/auth"
	rolesrp "go.microcore.dev/auth-service/internal/port/adapter/repository/roles"
	rulesrp "go.microcore.dev/auth-service/internal/port/adapter/repository/rules"
	devicessp "go.microcore.dev/auth-service/internal/port/service/devices"
	logoutsp "go.microcore.dev/auth-service/internal/port/service/logout"
	rolessp "go.microcore.dev/auth-service/internal/port/service/roles"
	rulessp "go.microcore.dev/auth-service/internal/port/service/rules"
	tokenssp "go.microcore.dev/auth-service/internal/port/service/tokens"
	devicessa "go.microcore.dev/auth-service/internal/service/devices"
	logoutsa "go.microcore.dev/auth-service/internal/service/logout"
	rolessa "go.microcore.dev/auth-service/internal/service/roles"
	rulessa "go.microcore.dev/auth-service/internal/service/rules"
	tokenssa "go.microcore.dev/auth-service/internal/service/tokens"
	sharedhttp "go.microcore.dev/auth-service/internal/shared/http/server"
	sharedpg "go.microcore.dev/auth-service/internal/shared/postgres"
	sharedredis "go.microcore.dev/auth-service/internal/shared/redis"
	sharedtel "go.microcore.dev/auth-service/internal/shared/telemetry"
	"go.microcore.dev/framework/log"
)

func Init(ctx context.Context, opts *Options) (*Server, error) {
	wire.Build(
		NewConfig,
		newLogger,
		newTelemetry,
		newPostgres,
		newRedis,
		newHTTPServer,
		// Repository
		newAuthRepository,
		newRolesRepository,
		newRulesRepository,
		// Services
		newTokensService,
		newRolesService,
		newRulesService,
		newDevicesService,
		newLogoutService,
		// Handlers
		newAuthHandler,
		newDevicesHandler,
		newLogoutHandler,
		newRolesHandler,
		newRulesHandler,
		newTokensHandler,

		newServer,
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

func newHTTPServer(
	telemetry *sharedtel.Telemetry,
) (*sharedhttp.Server, error) {
	s, err := sharedhttp.Init(
		&sharedhttp.Options{
			Telemetry: telemetry,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("init http server: %w", err)
	}
	s.Setup()
	return s, nil
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

func newDevicesService(
	authRepository authrp.Adapter,
) (devicessp.Service, error) {
	s, err := devicessa.Init(
		&devicessa.Options{
			AuthRepository: authRepository,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("init devices service: %w", err)
	}
	return s, nil
}

func newLogoutService(
	authRepository authrp.Adapter,
) (logoutsp.Service, error) {
	s, err := logoutsa.Init(
		&logoutsa.Options{
			AuthRepository: authRepository,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("init logout service: %w", err)
	}
	return s, nil
}

func newAuthHandler(
	tokensService tokenssp.Service,
	rolesService rolessp.Service,
) (authhp.Adapter, error) {
	h, err := authha.Init(
		&authha.Options{
			TokensService: tokensService,
			RolesService:  rolesService,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("init auth handler: %w", err)
	}
	return h, nil
}

func newDevicesHandler(
	devicesService devicessp.Service,
) (deviceshp.Adapter, error) {
	h, err := devicesha.Init(
		&devicesha.Options{
			DevicesService: devicesService,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("init devices handler: %w", err)
	}
	return h, nil
}

func newLogoutHandler(
	logoutService logoutsp.Service,
) (logouthp.Adapter, error) {
	h, err := logoutha.Init(
		&logoutha.Options{
			LogoutService: logoutService,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("init logout handler: %w", err)
	}
	return h, nil
}

func newRolesHandler(
	rolesService rolessp.Service,
) (roleshp.Adapter, error) {
	h, err := rolesha.Init(
		&rolesha.Options{
			RolesService: rolesService,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("init roles handler: %w", err)
	}
	return h, nil
}

func newRulesHandler(
	rulesService rulessp.Service,
) (ruleshp.Adapter, error) {
	h, err := rulesha.Init(
		&rulesha.Options{
			RulesService: rulesService,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("init rules handler: %w", err)
	}
	return h, nil
}

func newTokensHandler(
	tokensService tokenssp.Service,
	rolesService rolessp.Service,
) (tokenshp.Adapter, error) {
	h, err := tokensha.Init(
		&tokensha.Options{
			TokensService: tokensService,
			RolesService:  rolesService,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("init tokens handler: %w", err)
	}
	return h, nil
}

func newServer(
	config *Config,
	logger *slog.Logger,
	telemetry *sharedtel.Telemetry,
	postgres *sharedpg.Postgres,
	redis *sharedredis.Redis,
	httpServer *sharedhttp.Server,
	// Repository
	rolesRepository rolesrp.Adapter,
	// Handler
	authHandler authhp.Adapter,
	devicesHandler deviceshp.Adapter,
	logoutHandler logouthp.Adapter,
	rolesHandler roleshp.Adapter,
	rulesHandler ruleshp.Adapter,
	tokensHandler tokenshp.Adapter,
) *Server {
	return &Server{
		Config:     config,
		Logger:     logger,
		Telemetry:  telemetry,
		Postgres:   postgres,
		Redis:      redis,
		HTTPServer: httpServer,
		Repository: &Repository{
			Roles: rolesRepository,
		},
		Handler: &Handler{
			Auth:    authHandler,
			Devices: devicesHandler,
			Logout:  logoutHandler,
			Roles:   rolesHandler,
			Rules:   rulesHandler,
			Tokens:  tokensHandler,
		},
	}
}
