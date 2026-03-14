// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package bootstrap

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	authrp "go.microcore.dev/auth-service/internal/port/adapter/repository/auth"
	sharedtel "go.microcore.dev/auth-service/internal/shared/telemetry"
	"go.microcore.dev/framework/shutdown"
	"golang.org/x/term"
)

const (
	outputFileDirPerm = 0o755
	outputFilePerm    = 0o600
)

// ErrInvalidFlagValue is returned when a CLI flag has an invalid value.
var ErrInvalidFlagValue = errors.New("invalid flag value")

type (
	// Bootstrap represents a bootstrap application.
	Bootstrap struct {
		Options   *Options
		Config    *Config
		Logger    *slog.Logger
		Telemetry *sharedtel.Telemetry
		// Templates
		keyTpl             *template.Template
		securityWarningTpl *template.Template
	}

	// Options defines bootstrap primary options.
	Options struct {
		Output *OutputOptions
	}

	// OutputOptions configures output behavior for keys.
	OutputOptions struct {
		JSON          bool
		JSONPretty    bool
		Quiet         bool
		JwtAccessKey  string
		JwtRefreshKey string
		JwtHashKey    string
		AuthKey       string
		Stdout        io.Writer
	}

	// OutputConfig holds writers for various sensitive keys.
	OutputConfig struct {
		JwtAccessKeyWriter  io.Writer
		JwtRefreshKeyWriter io.Writer
		JwtHashKeyWriter    io.Writer
		AuthKeyWriter       io.Writer
	}

	// OutputData contains generated sensitive key values.
	OutputData struct {
		JwtAccessKey  string `json:"jwtAccessKey"`
		JwtRefreshKey string `json:"jwtRefreshKey"`
		JwtHashKey    string `json:"jwtHashKey"`
		AuthKey       string `json:"authKey"`
	}

	// OutputFlagWriterConfig holds configuration for key output destination.
	OutputFlagWriterConfig struct {
		Output string
		Field  string
		Flag   string
	}

	keyData struct {
		Sep         string
		Name        string
		Description string
		Value       string
	}

	securityWarningData struct {
		Sep string
	}
)

// Run executes the bootstrap application.
func (b *Bootstrap) Run(ctx context.Context) error {
	// Set key tpl
	b.SetKeyTpl()

	// Set securityWarningTpl
	b.SetSecurityWarningTpl()

	// Output config
	outputConfig, err := b.OutputConfig(ctx)
	if err != nil {
		return fmt.Errorf("output config: %w", err)
	}

	// Gen jwt access key
	jwtAccessKey, err := genKey(authrp.JWTSignKeyMinLen)
	if err != nil {
		return fmt.Errorf("gen jwt access key: %w", err)
	}

	// Gen jwt refresh key
	jwtRefreshKey, err := genKey(authrp.JWTSignKeyMinLen)
	if err != nil {
		return fmt.Errorf("gen jwt refresh key: %w", err)
	}

	// Gen jwt hash key
	jwtHashKey, err := genKey(authrp.JWTSignKeyMinLen)
	if err != nil {
		return fmt.Errorf("gen jwt hash key: %w", err)
	}

	// Gen auth key
	authKey, err := genKey(authrp.AuthKeyMinLen)
	if err != nil {
		return fmt.Errorf("gen auth key: %w", err)
	}

	// Output
	toB64Str := base64.StdEncoding.EncodeToString

	if err := b.Output(
		outputConfig,
		&OutputData{
			JwtAccessKey:  toB64Str(jwtAccessKey),
			JwtRefreshKey: toB64Str(jwtRefreshKey),
			JwtHashKey:    toB64Str(jwtHashKey),
			AuthKey:       toB64Str(authKey),
		},
	); err != nil {
		return fmt.Errorf("output: %w", err)
	}

	b.Logger.InfoContext(ctx, "bootstrap completed")

	return nil
}

// OutputConfig prepares writers for each key based on output options.
func (b *Bootstrap) OutputConfig(ctx context.Context) (*OutputConfig, error) {
	config := &OutputConfig{
		JwtAccessKeyWriter:  nil,
		JwtRefreshKeyWriter: nil,
		JwtHashKeyWriter:    nil,
		AuthKeyWriter:       nil,
	}

	if b.Options.Output.JSON {
		b.Logger.InfoContext(ctx, "outputting in JSON format to stdout")
		return config, nil
	}

	var err error

	// jwt access key
	if config.JwtAccessKeyWriter, err = b.OutputFlagWriter(
		ctx,
		&OutputFlagWriterConfig{
			Output: b.Options.Output.JwtAccessKey,
			Field:  "jwt access key",
			Flag:   "jwt-access-key-output",
		},
	); err != nil {
		return nil, fmt.Errorf("jwt access key writer : %w", err)
	}

	// jwt refresh key
	if config.JwtRefreshKeyWriter, err = b.OutputFlagWriter(
		ctx,
		&OutputFlagWriterConfig{
			Output: b.Options.Output.JwtRefreshKey,
			Field:  "jwt refresh key",
			Flag:   "jwt-refresh-key-output",
		},
	); err != nil {
		return nil, fmt.Errorf("jwt refresh key writer : %w", err)
	}

	// jwt hash key
	if config.JwtHashKeyWriter, err = b.OutputFlagWriter(
		ctx,
		&OutputFlagWriterConfig{
			Output: b.Options.Output.JwtHashKey,
			Field:  "jwt hash key",
			Flag:   "jwt-hash-key-output",
		},
	); err != nil {
		return nil, fmt.Errorf("jwt hash key writer : %w", err)
	}

	// auth key
	if config.AuthKeyWriter, err = b.OutputFlagWriter(
		ctx,
		&OutputFlagWriterConfig{
			Output: b.Options.Output.AuthKey,
			Field:  "auth key",
			Flag:   "auth-key-output",
		},
	); err != nil {
		return nil, fmt.Errorf("auth key writer : %w", err)
	}

	return config, nil
}

// OutputFlagWriter returns an [io.Writer] for a key, either stdout or file.
func (b *Bootstrap) OutputFlagWriter(
	ctx context.Context,
	config *OutputFlagWriterConfig,
) (io.Writer, error) {
	switch {
	case config.Output == "" || config.Output == "stdout":
		b.Logger.InfoContext(
			ctx,
			config.Field,
			slog.String("output", "stdout"),
			slog.Bool("quiet", b.Options.Output.Quiet),
		)

		return b.Options.Output.Stdout, nil

	case strings.HasPrefix(config.Output, "file:"):
		return b.openOutputFile(ctx, config)
	default:
		return nil, fmt.Errorf("%w: --%s=%s", ErrInvalidFlagValue, config.Flag, config.Output)
	}
}

// Output writes keys to the configured writers or JSON output.
func (b *Bootstrap) Output(config *OutputConfig, data *OutputData) error {
	switch {
	case b.Options.Output.JSON:
		return b.outputJSONHandler(data)
	case b.Options.Output.Quiet:
		return b.outputQuietHandler(config, data)
	default:
		return b.outputDefaultHandler(config, data)
	}
}

// OutputQuietWrite prints a value quietly to the given writer.
func (b *Bootstrap) OutputQuietWrite(w io.Writer, s string) error {
	fn := fmt.Fprintln
	if w != b.Options.Output.Stdout {
		fn = fmt.Fprint
	}

	if _, err := fn(w, s); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}

// SetKeyTpl initializes the template used to render key output
// with name, description, value, and separator formatting.
func (b *Bootstrap) SetKeyTpl() {
	b.keyTpl = template.Must(template.New("keyOutput").Parse(`{{.Sep}}
{{.Name}}

{{.Description}}

{{.Value}}
{{.Sep}}`))
}

// SetSecurityWarningTpl initializes the template used to display
// a one-time security warning for sensitive output.
func (b *Bootstrap) SetSecurityWarningTpl() {
	b.securityWarningTpl = template.Must(template.New("securityWarningOutput").Parse(
		`{{.Sep}}
SECURITY WARNING!

The following output contains highly sensitive information intended for the operator only. ` +
			`This data is displayed ONCE and must be stored securely. Do NOT share it, commit it to ` +
			`version control, or expose it in logs, CI systems, or monitoring tools.
{{.Sep}}`,
	))
}

func (b *Bootstrap) outputJSONHandler(data *OutputData) error {
	if err := b.outputJSON(data); err != nil {
		return fmt.Errorf("output json: %w", err)
	}

	return nil
}

func (b *Bootstrap) outputQuietHandler(config *OutputConfig, data *OutputData) error {
	if err := b.outputQuietPrint(config, data); err != nil {
		return fmt.Errorf("output quiet print: %w", err)
	}

	return nil
}

func (b *Bootstrap) outputDefaultHandler(config *OutputConfig, data *OutputData) error {
	if err := b.outputSecurityWarning(); err != nil {
		return fmt.Errorf("output security warning: %w", err)
	}

	if err := b.outputPrint(config, data); err != nil {
		return fmt.Errorf("output print: %w", err)
	}

	return nil
}

func (b *Bootstrap) outputJSON(data *OutputData) error {
	enc := json.NewEncoder(b.Options.Output.Stdout)
	if b.Options.Output.JSONPretty {
		enc.SetIndent("", "  ")
	}

	//nolint:gosec // AuthKey is a public identifier/config, not a raw private secret
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("marshal output data: %w", err)
	}

	return nil
}

func (b *Bootstrap) outputQuietPrint(config *OutputConfig, data *OutputData) error {
	if err := b.OutputQuietWrite(config.JwtAccessKeyWriter, data.JwtAccessKey); err != nil {
		return fmt.Errorf("output quiet print: jwt access key: %w", err)
	}

	if err := b.OutputQuietWrite(config.JwtRefreshKeyWriter, data.JwtRefreshKey); err != nil {
		return fmt.Errorf("output quiet print: jwt refresh key: %w", err)
	}

	if err := b.OutputQuietWrite(config.JwtHashKeyWriter, data.JwtHashKey); err != nil {
		return fmt.Errorf("output quiet print: jwt hash key: %w", err)
	}

	if err := b.OutputQuietWrite(config.AuthKeyWriter, data.AuthKey); err != nil {
		return fmt.Errorf("output quiet print: auth key: %w", err)
	}

	return nil
}

func (b *Bootstrap) outputPrint(config *OutputConfig, data *OutputData) error {
	if err := b.outputKey(
		config.JwtAccessKeyWriter,
		"JWT_ACCESS_KEY",
		"This key is used to generate and validate short-lived access tokens required "+
			"for API authentication and authorization. Compromise allows issuing valid tokens.",
		data.JwtAccessKey,
	); err != nil {
		return fmt.Errorf("output key: jwt access key: %w", err)
	}

	if err := b.outputKey(
		config.JwtRefreshKeyWriter,
		"JWT_REFRESH_KEY",
		"This key is used to sign refresh tokens allowing access token renewal without "+
			"re-authentication. Compromise allows unlimited session renewal.",
		data.JwtRefreshKey,
	); err != nil {
		return fmt.Errorf("output key: jwt refresh key: %w", err)
	}

	if err := b.outputKey(
		config.JwtHashKeyWriter,
		"JWT_HASH_KEY",
		"This key protects permanent tokens by encrypting them before storage. Compromise "+
			"allows decryption of all long-lived tokens.",
		data.JwtHashKey,
	); err != nil {
		return fmt.Errorf("output key: jwt hash key: %w", err)
	}

	if err := b.outputKey(
		config.AuthKeyWriter,
		"AUTH_KEY",
		"This key secures inter-service requests when requesting token issuance from the "+
			"Auth Service. Compromise allows unauthorized requests.",
		data.AuthKey,
	); err != nil {
		return fmt.Errorf("output key: auth key: %w", err)
	}

	return nil
}

func (b *Bootstrap) outputSecurityWarning() error {
	if err := b.securityWarningTpl.Execute(
		b.Options.Output.Stdout,
		securityWarningData{
			Sep: getSep(),
		},
	); err != nil {
		return fmt.Errorf("template execute: %w", err)
	}

	return nil
}

func (b *Bootstrap) openOutputFile(
	ctx context.Context,
	config *OutputFlagWriterConfig,
) (io.Writer, error) {
	path := strings.TrimPrefix(config.Output, "file:")
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, outputFileDirPerm); err != nil {
		return nil, fmt.Errorf("create output directory: %w", err)
	}

	f, err := os.OpenFile(
		filepath.Clean(path),
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		outputFilePerm,
	)
	if err != nil {
		return nil, fmt.Errorf("create output file: %w", err)
	}

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

	b.Logger.InfoContext(
		ctx,
		config.Field,
		slog.String("output", "file"),
		slog.String("path", path),
		slog.Bool("quiet", b.Options.Output.Quiet),
	)

	return f, nil
}

func (b *Bootstrap) outputKey(writer io.Writer, name, description, value string) error {
	if err := b.keyTpl.Execute(
		writer,
		keyData{
			Sep:         getSep(),
			Name:        name,
			Description: description,
			Value:       value,
		},
	); err != nil {
		return fmt.Errorf("template execute: %w", err)
	}

	return nil
}

// genKey generates a cryptographically secure random key of the given size.
func genKey(size int) ([]byte, error) {
	key := make([]byte, size)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("crypto/rand failed: %w", err)
	}

	return key, nil
}

// getSep returns the string separator.
func getSep() string {
	return strings.Repeat("-", termWidth())
}

// termWidth returns the terminal width or defaults to 80 columns.
func termWidth() int {
	defaultWidth := 80

	fd := os.Stderr.Fd()
	if uint64(fd) > uint64(^uint(0)>>1) {
		return defaultWidth
	}

	width, _, err := term.GetSize(int(fd))
	if err != nil || width <= 0 {
		width = defaultWidth
	}

	return width
}
