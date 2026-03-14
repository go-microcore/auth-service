// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package bootstrap_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.microcore.dev/auth-service/internal/app/bootstrap"
)

const (
	jwtAccessKey  = "jwt.access.key"
	jwtRefreshKey = "jwt.refresh.key"
	jwtHashKey    = "jwt.hash.key"
	authKey       = "auth.key"
)

func Test_OutputSettings_JSON(t *testing.T) {
	t.Parallel()

	app := &bootstrap.Bootstrap{
		Options: &bootstrap.Options{
			Output: &bootstrap.OutputOptions{
				JSON:          true,
				JSONPretty:    false,
				Quiet:         false,
				JwtAccessKey:  "key",
				JwtRefreshKey: "key",
				JwtHashKey:    "key",
				AuthKey:       "key",
				Stdout:        nil,
			},
		},
		Config: &bootstrap.Config{
			Name: "name",
		},
		Logger:    slog.New(slog.DiscardHandler),
		Telemetry: nil,
	}

	app.SetKeyTpl()
	app.SetSecurityWarningTpl()

	settings, err := app.OutputConfig(t.Context())
	require.NoError(t, err)
	require.NotNil(t, settings)
	require.Nil(t, settings.JwtAccessKeyWriter)
	require.Nil(t, settings.JwtRefreshKeyWriter)
	require.Nil(t, settings.JwtHashKeyWriter)
	require.Nil(t, settings.AuthKeyWriter)
}

func Test_OutputFlagWriter(t *testing.T) {
	t.Parallel()

	t.Run("blank", func(t *testing.T) {
		t.Parallel()

		app := &bootstrap.Bootstrap{
			Options: &bootstrap.Options{
				Output: &bootstrap.OutputOptions{
					JSON:          false,
					JSONPretty:    false,
					Quiet:         false,
					JwtAccessKey:  "key",
					JwtRefreshKey: "key",
					JwtHashKey:    "key",
					AuthKey:       "key",
					Stdout:        os.Stdout,
				},
			},
			Config: &bootstrap.Config{
				Name: "name",
			},
			Logger:    slog.New(slog.DiscardHandler),
			Telemetry: nil,
		}

		app.SetKeyTpl()
		app.SetSecurityWarningTpl()

		writer, err := app.OutputFlagWriter(
			t.Context(),
			&bootstrap.OutputFlagWriterConfig{
				Output: "",
				Field:  "field",
				Flag:   "flag",
			},
		)
		require.NoError(t, err)
		require.NotNil(t, writer)
		require.Equal(t, os.Stdout, writer)
	})

	t.Run("stdout", func(t *testing.T) {
		t.Parallel()

		app := &bootstrap.Bootstrap{
			Options: &bootstrap.Options{
				Output: &bootstrap.OutputOptions{
					JSON:          false,
					JSONPretty:    false,
					Quiet:         false,
					JwtAccessKey:  "key",
					JwtRefreshKey: "key",
					JwtHashKey:    "key",
					AuthKey:       "key",
					Stdout:        os.Stdout,
				},
			},
			Config: &bootstrap.Config{
				Name: "name",
			},
			Logger:    slog.New(slog.DiscardHandler),
			Telemetry: nil,
		}

		app.SetKeyTpl()
		app.SetSecurityWarningTpl()

		writer, err := app.OutputFlagWriter(
			t.Context(),
			&bootstrap.OutputFlagWriterConfig{
				Output: "stdout",
				Field:  "field",
				Flag:   "flag",
			},
		)
		require.NoError(t, err)
		require.NotNil(t, writer)
		require.Equal(t, os.Stdout, writer)
	})

	t.Run("file", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "token.txt")

		app := &bootstrap.Bootstrap{
			Options: &bootstrap.Options{
				Output: &bootstrap.OutputOptions{
					JSON:          false,
					JSONPretty:    false,
					Quiet:         false,
					JwtAccessKey:  "key",
					JwtRefreshKey: "key",
					JwtHashKey:    "key",
					AuthKey:       "key",
					Stdout:        os.Stdout,
				},
			},
			Config: &bootstrap.Config{
				Name: "name",
			},
			Logger:    slog.New(slog.DiscardHandler),
			Telemetry: nil,
		}

		app.SetKeyTpl()
		app.SetSecurityWarningTpl()

		writer, err := app.OutputFlagWriter(
			t.Context(),
			&bootstrap.OutputFlagWriterConfig{
				Output: "file:" + path,
				Field:  "field",
				Flag:   "flag",
			},
		)
		require.NoError(t, err)
		require.NotNil(t, writer)

		_, ok := writer.(*os.File)
		require.True(t, ok)

		_, err = os.Stat(path)
		require.NoError(t, err)

		fi, err := os.Stat(path)
		require.NoError(t, err)
		assert.Equal(t, os.FileMode(0o600), fi.Mode().Perm())
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		app := &bootstrap.Bootstrap{
			Options: &bootstrap.Options{
				Output: &bootstrap.OutputOptions{
					JSON:          false,
					JSONPretty:    false,
					Quiet:         false,
					JwtAccessKey:  "key",
					JwtRefreshKey: "key",
					JwtHashKey:    "key",
					AuthKey:       "key",
					Stdout:        os.Stdout,
				},
			},
			Config: &bootstrap.Config{
				Name: "name",
			},
			Logger:    slog.New(slog.DiscardHandler),
			Telemetry: nil,
		}

		app.SetKeyTpl()
		app.SetSecurityWarningTpl()

		writer, err := app.OutputFlagWriter(
			t.Context(),
			&bootstrap.OutputFlagWriterConfig{
				Output: "invalid",
				Field:  "field",
				Flag:   "flag",
			},
		)
		require.Error(t, err)
		require.Nil(t, writer)
	})
}

func Test_Output_JSON(t *testing.T) {
	t.Parallel()

	t.Run("compact", func(t *testing.T) {
		t.Parallel()

		buf := &bytes.Buffer{}

		app := &bootstrap.Bootstrap{
			Options: &bootstrap.Options{
				Output: &bootstrap.OutputOptions{
					JSON:          true,
					JSONPretty:    false,
					Quiet:         false,
					JwtAccessKey:  "key",
					JwtRefreshKey: "key",
					JwtHashKey:    "key",
					AuthKey:       "key",
					Stdout:        buf,
				},
			},
			Config: &bootstrap.Config{
				Name: "name",
			},
			Logger:    nil,
			Telemetry: nil,
		}

		app.SetKeyTpl()
		app.SetSecurityWarningTpl()

		err := app.Output(nil, &bootstrap.OutputData{
			JwtAccessKey:  jwtAccessKey,
			JwtRefreshKey: jwtRefreshKey,
			JwtHashKey:    jwtHashKey,
			AuthKey:       authKey,
		})
		require.NoError(t, err)

		var got map[string]string

		err = json.Unmarshal(buf.Bytes(), &got)
		require.NoError(t, err)

		assert.Equal(t, jwtAccessKey, got["jwtAccessKey"])
		assert.Equal(t, jwtRefreshKey, got["jwtRefreshKey"])
		assert.Equal(t, jwtHashKey, got["jwtHashKey"])
		assert.Equal(t, authKey, got["authKey"])
	})

	t.Run("pretty", func(t *testing.T) {
		t.Parallel()

		buf := &bytes.Buffer{}

		app := &bootstrap.Bootstrap{
			Options: &bootstrap.Options{
				Output: &bootstrap.OutputOptions{
					JSON:          true,
					JSONPretty:    true,
					Quiet:         false,
					JwtAccessKey:  "key",
					JwtRefreshKey: "key",
					JwtHashKey:    "key",
					AuthKey:       "key",
					Stdout:        buf,
				},
			},
			Config: &bootstrap.Config{
				Name: "name",
			},
			Logger:    nil,
			Telemetry: nil,
		}

		app.SetKeyTpl()
		app.SetSecurityWarningTpl()

		err := app.Output(nil, &bootstrap.OutputData{
			JwtAccessKey:  jwtAccessKey,
			JwtRefreshKey: jwtRefreshKey,
			JwtHashKey:    jwtHashKey,
			AuthKey:       authKey,
		})
		require.NoError(t, err)

		var got map[string]string

		err = json.Unmarshal(buf.Bytes(), &got)
		require.NoError(t, err)

		assert.Equal(t, jwtAccessKey, got["jwtAccessKey"])
		assert.Equal(t, jwtRefreshKey, got["jwtRefreshKey"])
		assert.Equal(t, jwtHashKey, got["jwtHashKey"])
		assert.Equal(t, authKey, got["authKey"])
	})
}

func Test_Output_Custom(t *testing.T) {
	t.Parallel()

	t.Run("default", func(t *testing.T) {
		t.Parallel()

		jwtAccessKeyBuf := &bytes.Buffer{}
		jwtRefreshKeyBuf := &bytes.Buffer{}
		jwtHashKeyBuf := &bytes.Buffer{}
		authKeyBuf := &bytes.Buffer{}

		app := &bootstrap.Bootstrap{
			Options: &bootstrap.Options{
				Output: &bootstrap.OutputOptions{
					JSON:          false,
					JSONPretty:    false,
					Quiet:         false,
					JwtAccessKey:  "key",
					JwtRefreshKey: "key",
					JwtHashKey:    "key",
					AuthKey:       "key",
					Stdout:        os.Stdout,
				},
			},
			Config: &bootstrap.Config{
				Name: "name",
			},
			Logger:    nil,
			Telemetry: nil,
		}

		app.SetKeyTpl()
		app.SetSecurityWarningTpl()

		err := app.Output(
			&bootstrap.OutputConfig{
				JwtAccessKeyWriter:  jwtAccessKeyBuf,
				JwtRefreshKeyWriter: jwtRefreshKeyBuf,
				JwtHashKeyWriter:    jwtHashKeyBuf,
				AuthKeyWriter:       authKeyBuf,
			},
			&bootstrap.OutputData{
				JwtAccessKey:  jwtAccessKey,
				JwtRefreshKey: jwtRefreshKey,
				JwtHashKey:    jwtHashKey,
				AuthKey:       authKey,
			},
		)
		require.NoError(t, err)

		assert.Contains(t, jwtAccessKeyBuf.String(), jwtAccessKey)
		assert.Contains(t, jwtRefreshKeyBuf.String(), jwtRefreshKey)
		assert.Contains(t, jwtHashKeyBuf.String(), jwtHashKey)
		assert.Contains(t, authKeyBuf.String(), authKey)
	})

	t.Run("quiet", func(t *testing.T) {
		t.Parallel()

		jwtAccessKeyBuf := &bytes.Buffer{}
		jwtRefreshKeyBuf := &bytes.Buffer{}
		jwtHashKeyBuf := &bytes.Buffer{}
		authKeyBuf := &bytes.Buffer{}

		app := &bootstrap.Bootstrap{
			Options: &bootstrap.Options{
				Output: &bootstrap.OutputOptions{
					JSON:          false,
					JSONPretty:    false,
					Quiet:         true,
					JwtAccessKey:  "key",
					JwtRefreshKey: "key",
					JwtHashKey:    "key",
					AuthKey:       "key",
					Stdout:        os.Stdout,
				},
			},
			Config: &bootstrap.Config{
				Name: "name",
			},
			Logger:    nil,
			Telemetry: nil,
		}

		app.SetKeyTpl()
		app.SetSecurityWarningTpl()

		err := app.Output(
			&bootstrap.OutputConfig{
				JwtAccessKeyWriter:  jwtAccessKeyBuf,
				JwtRefreshKeyWriter: jwtRefreshKeyBuf,
				JwtHashKeyWriter:    jwtHashKeyBuf,
				AuthKeyWriter:       authKeyBuf,
			},
			&bootstrap.OutputData{
				JwtAccessKey:  jwtAccessKey,
				JwtRefreshKey: jwtRefreshKey,
				JwtHashKey:    jwtHashKey,
				AuthKey:       authKey,
			},
		)
		require.NoError(t, err)

		assert.Equal(t, jwtAccessKey, jwtAccessKeyBuf.String())
		assert.Equal(t, jwtRefreshKey, jwtRefreshKeyBuf.String())
		assert.Equal(t, jwtHashKey, jwtHashKeyBuf.String())
		assert.Equal(t, authKey, authKeyBuf.String())
	})
}
