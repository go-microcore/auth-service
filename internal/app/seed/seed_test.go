// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package seed_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.microcore.dev/auth-service/internal/app/seed"
	rolessp "go.microcore.dev/auth-service/internal/port/service/roles"
	rulessp "go.microcore.dev/auth-service/internal/port/service/rules"
	tokenssp "go.microcore.dev/auth-service/internal/port/service/tokens"
	sharedhttp "go.microcore.dev/auth-service/internal/shared/http/client"
	sharedpg "go.microcore.dev/auth-service/internal/shared/postgres"
	sharedredis "go.microcore.dev/auth-service/internal/shared/redis"
	sharedtel "go.microcore.dev/auth-service/internal/shared/telemetry"
)

const (
	adminStaticToken = "fake.token"
)

type (
	SeedBuilder struct {
		t *testing.T

		Options    *seed.Options
		Config     *seed.Config
		Logger     *slog.Logger
		Telemetry  *sharedtel.Telemetry
		Postgres   *sharedpg.Postgres
		Redis      *sharedredis.Redis
		HTTPClient *sharedhttp.Client
		Service    *seed.Service
	}
)

func Test_CreateAdminStaticToken_HappyPath(t *testing.T) {
	t.Parallel()

	fakeToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
	mockTokens := tokenssp.NewMockService(t)

	mockTokens.EXPECT().
		FilterStaticAccessTokens(mock.Anything, mock.Anything).
		Return([]tokenssp.StaticAccessTokenResult{}, nil)

	mockTokens.EXPECT().
		CreateStaticAccessToken(mock.Anything, mock.Anything).
		Return(fakeToken, nil)

	app := NewSeedBuilder(t).
		WithConfig(nil).
		WithService(&seed.Service{
			Tokens: mockTokens,
			Roles:  nil,
			Rules:  nil,
		}).
		WithLogger(slog.New(slog.DiscardHandler)).
		Build()

	token, err := app.CreateAdminStaticToken(t.Context())
	require.NoError(t, err)
	require.Equal(t, fakeToken, token)
}

func Test_CreateAdminStaticToken_AlreadyExists(t *testing.T) {
	t.Parallel()

	mockTokens := tokenssp.NewMockService(t)

	mockTokens.EXPECT().
		FilterStaticAccessTokens(mock.Anything, mock.Anything).
		Return([]tokenssp.StaticAccessTokenResult{{}}, nil)

	app := NewSeedBuilder(t).
		WithConfig(nil).
		WithService(&seed.Service{
			Tokens: mockTokens,
			Roles:  nil,
			Rules:  nil,
		}).
		Build()

	token, err := app.CreateAdminStaticToken(t.Context())
	require.NoError(t, err)
	require.Empty(t, token)
}

func Test_CreateAdminRole_HappyPath(t *testing.T) {
	t.Parallel()

	mockRoles := rolessp.NewMockService(t)

	mockRoles.EXPECT().
		FilterRoles(mock.Anything, mock.Anything).
		Return([]rolessp.FilterRolesResult{}, nil)

	mockRoles.EXPECT().
		CreateRole(mock.Anything, mock.Anything).
		Return(&rolessp.CreateRoleResult{}, nil)

	app := NewSeedBuilder(t).
		WithConfig(nil).
		WithService(&seed.Service{
			Tokens: nil,
			Roles:  mockRoles,
			Rules:  nil,
		}).
		WithLogger(slog.New(slog.DiscardHandler)).
		Build()

	err := app.CreateAdminRole(t.Context())
	require.NoError(t, err)
}

func Test_CreateAdminRole_AlreadyExists(t *testing.T) {
	t.Parallel()

	mockRoles := rolessp.NewMockService(t)

	mockRoles.EXPECT().
		FilterRoles(mock.Anything, mock.Anything).
		Return([]rolessp.FilterRolesResult{{}}, nil)

	app := NewSeedBuilder(t).
		WithConfig(nil).
		WithService(&seed.Service{
			Tokens: nil,
			Roles:  mockRoles,
			Rules:  nil,
		}).
		Build()

	err := app.CreateAdminRole(t.Context())
	require.NoError(t, err)
}

func Test_CreateAdminHTTPRule_HappyPath(t *testing.T) {
	t.Parallel()

	mockRules := rulessp.NewMockService(t)

	mockRules.EXPECT().
		FilterHTTPRules(mock.Anything, mock.Anything).
		Return([]rulessp.FilterHTTPRulesResult{}, nil)

	mockRules.EXPECT().
		CreateHTTPRule(mock.Anything, mock.Anything).
		Return(&rulessp.CreateHTTPRuleResult{}, nil)

	app := NewSeedBuilder(t).
		WithConfig(nil).
		WithService(&seed.Service{
			Tokens: nil,
			Rules:  mockRules,
			Roles:  nil,
		}).
		WithLogger(slog.New(slog.DiscardHandler)).
		Build()

	err := app.CreateAdminHTTPRule(t.Context())
	require.NoError(t, err)
}

func Test_CreateAdminHTTPRule_AlreadyExists(t *testing.T) {
	t.Parallel()

	mockRules := rulessp.NewMockService(t)

	mockRules.EXPECT().
		FilterHTTPRules(mock.Anything, mock.Anything).
		Return([]rulessp.FilterHTTPRulesResult{{}}, nil)

	app := NewSeedBuilder(t).
		WithConfig(nil).
		WithService(&seed.Service{
			Tokens: nil,
			Rules:  mockRules,
			Roles:  nil,
		}).
		Build()

	err := app.CreateAdminHTTPRule(t.Context())
	require.NoError(t, err)
}

func Test_OutputConfig_JSON(t *testing.T) {
	t.Parallel()

	app := NewSeedBuilder(t).
		WithConfig(nil).
		WithOptions(&seed.Options{
			Output: &seed.OutputOptions{
				JSON:       true,
				JSONPretty: false,
				Quiet:      false,
				AdminToken: "",
				Stdout:     nil,
			},
		}).
		WithLogger(slog.New(slog.DiscardHandler)).
		Build()

	config, err := app.OutputConfig(t.Context())
	require.NoError(t, err)
	require.Nil(t, config.AdminTokenWriter)
}

func Test_OutputConfig_AdminToken(t *testing.T) {
	t.Parallel()

	t.Run("blank", func(t *testing.T) {
		t.Parallel()

		app := NewSeedBuilder(t).
			WithConfig(nil).
			WithOptions(&seed.Options{
				Output: &seed.OutputOptions{
					JSON:       false,
					JSONPretty: false,
					Quiet:      false,
					AdminToken: "",
					Stdout:     os.Stdout,
				},
			}).
			WithLogger(slog.New(slog.DiscardHandler)).
			Build()

		config, err := app.OutputConfig(t.Context())
		require.NoError(t, err)
		require.NotNil(t, config)
		require.Equal(t, os.Stdout, config.AdminTokenWriter)
	})

	t.Run("stdout", func(t *testing.T) {
		t.Parallel()

		app := NewSeedBuilder(t).
			WithConfig(nil).
			WithOptions(&seed.Options{
				Output: &seed.OutputOptions{
					JSON:       false,
					JSONPretty: false,
					Quiet:      false,
					AdminToken: "stdout",
					Stdout:     os.Stdout,
				},
			}).
			WithLogger(slog.New(slog.DiscardHandler)).
			Build()

		config, err := app.OutputConfig(t.Context())
		require.NoError(t, err)
		require.NotNil(t, config)
		require.Equal(t, os.Stdout, config.AdminTokenWriter)
	})

	t.Run("file", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "token.txt")

		app := NewSeedBuilder(t).
			WithConfig(nil).
			WithOptions(&seed.Options{
				Output: &seed.OutputOptions{
					JSON:       false,
					JSONPretty: false,
					Quiet:      false,
					AdminToken: "file:" + path,
					Stdout:     nil,
				},
			}).
			WithLogger(slog.New(slog.DiscardHandler)).
			Build()

		config, err := app.OutputConfig(t.Context())
		require.NoError(t, err)
		require.NotNil(t, config)
		require.NotNil(t, config.AdminTokenWriter)

		_, ok := config.AdminTokenWriter.(*os.File)
		require.True(t, ok)

		_, err = os.Stat(path)
		require.NoError(t, err)

		fi, err := os.Stat(path)
		require.NoError(t, err)
		assert.Equal(t, os.FileMode(0o600), fi.Mode().Perm())
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		app := NewSeedBuilder(t).
			WithConfig(nil).
			WithOptions(&seed.Options{
				Output: &seed.OutputOptions{
					JSON:       false,
					JSONPretty: false,
					Quiet:      false,
					AdminToken: "invalid",
					Stdout:     nil,
				},
			}).
			WithLogger(slog.New(slog.DiscardHandler)).
			Build()

		config, err := app.OutputConfig(t.Context())
		require.Error(t, err)
		require.Nil(t, config)
	})
}

func Test_Output_JSON(t *testing.T) {
	t.Parallel()

	t.Run("compact", func(t *testing.T) {
		t.Parallel()

		buf := &bytes.Buffer{}

		app := NewSeedBuilder(t).
			WithConfig(nil).
			WithOptions(&seed.Options{
				Output: &seed.OutputOptions{
					JSON:       true,
					JSONPretty: false,
					Quiet:      false,
					AdminToken: "",
					Stdout:     buf,
				},
			}).
			Build()

		err := app.Output(nil, &seed.OutputData{AdminStaticToken: "fake.token"})
		require.NoError(t, err)

		var got map[string]string

		err = json.Unmarshal(buf.Bytes(), &got)
		require.NoError(t, err)

		assert.Equal(t, "fake.token", got["adminToken"])
	})

	t.Run("pretty", func(t *testing.T) {
		t.Parallel()

		buf := &bytes.Buffer{}

		app := NewSeedBuilder(t).
			WithConfig(nil).
			WithOptions(&seed.Options{
				Output: &seed.OutputOptions{
					JSON:       true,
					JSONPretty: true,
					Quiet:      false,
					AdminToken: "",
					Stdout:     buf,
				},
			}).
			Build()

		err := app.Output(nil, &seed.OutputData{AdminStaticToken: "fake.token"})
		require.NoError(t, err)

		var got map[string]string

		err = json.Unmarshal(buf.Bytes(), &got)
		require.NoError(t, err)

		assert.Equal(t, "fake.token", got["adminToken"])
	})
}

func Test_Output_AdminTokenWriter(t *testing.T) {
	t.Parallel()

	t.Run("default", func(t *testing.T) {
		t.Parallel()

		buf := &bytes.Buffer{}

		app := NewSeedBuilder(t).
			WithConfig(nil).
			WithOptions(&seed.Options{
				Output: &seed.OutputOptions{
					JSON:       false,
					JSONPretty: false,
					Quiet:      false,
					AdminToken: "",
					Stdout:     nil,
				},
			}).
			Build()

		adminStaticToken := "fake.token"

		err := app.Output(
			&seed.OutputConfig{AdminTokenWriter: buf},
			&seed.OutputData{AdminStaticToken: adminStaticToken},
		)
		require.NoError(t, err)

		assert.Contains(t, buf.String(), adminStaticToken)
	})

	t.Run("quiet", func(t *testing.T) {
		t.Parallel()

		buf := &bytes.Buffer{}

		app := NewSeedBuilder(t).
			WithConfig(nil).
			WithOptions(&seed.Options{
				Output: &seed.OutputOptions{
					JSON:       false,
					JSONPretty: false,
					Quiet:      true,
					AdminToken: "",
					Stdout:     nil,
				},
			}).
			Build()

		err := app.Output(
			&seed.OutputConfig{AdminTokenWriter: buf},
			&seed.OutputData{AdminStaticToken: adminStaticToken},
		)
		require.NoError(t, err)

		assert.Equal(t, adminStaticToken, buf.String())
	})
}

func NewSeedBuilder(t *testing.T) *SeedBuilder {
	t.Helper()

	return &SeedBuilder{
		t:          t,
		Options:    nil,
		Config:     nil,
		Logger:     nil,
		Telemetry:  nil,
		Postgres:   nil,
		Redis:      nil,
		HTTPClient: nil,
		Service:    nil,
	}
}

func (b *SeedBuilder) WithOptions(o *seed.Options) *SeedBuilder {
	b.Options = o
	return b
}

func (b *SeedBuilder) WithConfig(c *seed.Config) *SeedBuilder {
	if c == nil {
		c = &seed.Config{
			Name: "name",
			Admin: &seed.ConfigAdmin{
				Role: &seed.ConfigAdminRole{
					ID:          "id",
					Name:        "name",
					Description: "description",
				},
				Rule: &seed.ConfigAdminRule{
					Path:    "id",
					Methods: []string{"name"},
				},
				Token: &seed.ConfigAdminToken{
					ID:          "id",
					Description: "description",
				},
			},
		}
	}

	b.Config = c

	return b
}

func (b *SeedBuilder) WithLogger(l *slog.Logger) *SeedBuilder {
	b.Logger = l
	return b
}

func (b *SeedBuilder) WithTelemetry(t *sharedtel.Telemetry) *SeedBuilder {
	b.Telemetry = t
	return b
}

func (b *SeedBuilder) WithPostgres(p *sharedpg.Postgres) *SeedBuilder {
	b.Postgres = p
	return b
}

func (b *SeedBuilder) WithRedis(r *sharedredis.Redis) *SeedBuilder {
	b.Redis = r
	return b
}

func (b *SeedBuilder) WithHTTPClient(c *sharedhttp.Client) *SeedBuilder {
	b.HTTPClient = c
	return b
}

func (b *SeedBuilder) WithService(s *seed.Service) *SeedBuilder {
	b.Service = s
	return b
}

func (b *SeedBuilder) Build() *seed.Seed {
	app := &seed.Seed{
		Options:    b.Options,
		Config:     b.Config,
		Logger:     b.Logger,
		Telemetry:  b.Telemetry,
		Postgres:   b.Postgres,
		Redis:      b.Redis,
		HTTPClient: b.HTTPClient,
		Service:    b.Service,
	}

	app.SetOutputTokenTpl()

	return app
}
