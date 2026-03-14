// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package seed

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	rolessp "go.microcore.dev/auth-service/internal/port/service/roles"
	rulessp "go.microcore.dev/auth-service/internal/port/service/rules"
	tokenssp "go.microcore.dev/auth-service/internal/port/service/tokens"
	sharedhttp "go.microcore.dev/auth-service/internal/shared/http/client"
	sharedpg "go.microcore.dev/auth-service/internal/shared/postgres"
	sharedredis "go.microcore.dev/auth-service/internal/shared/redis"
	sharedtel "go.microcore.dev/auth-service/internal/shared/telemetry"
	"go.microcore.dev/framework/shutdown"
	"golang.org/x/term"
)

const (
	outputFileDirPerm  = 0o755
	outputFilePerm     = 0o600
	noStaticAdminToken = "(no static admin token created)"
)

// ErrInvalidAdminTokenOutput signals that the provided admin token output option is invalid.
var ErrInvalidAdminTokenOutput = errors.New("invalid admin token output value")

type (
	// Seed represents a seed application.
	Seed struct {
		Options    *Options
		Config     *Config
		Logger     *slog.Logger
		Telemetry  *sharedtel.Telemetry
		Postgres   *sharedpg.Postgres
		Redis      *sharedredis.Redis
		HTTPClient *sharedhttp.Client
		Service    *Service
		// Templates
		outputTokenTpl *template.Template
	}

	// Service wraps services.
	Service struct {
		Tokens tokenssp.Service
		Roles  rolessp.Service
		Rules  rulessp.Service
	}

	// Options defines seed primary options.
	Options struct {
		Output *OutputOptions
	}

	// OutputOptions configures output behavior for the seed.
	OutputOptions struct {
		JSON       bool
		JSONPretty bool
		Quiet      bool
		AdminToken string
		Stdout     io.Writer
	}

	// OutputConfig holds writers for outputting sensitive data.
	OutputConfig struct {
		AdminTokenWriter io.Writer
	}

	// OutputData contains generated output data.
	OutputData struct {
		AdminStaticToken string `json:"adminToken"`
	}

	outputTokenData struct {
		Sep   string
		Token string
	}
)

// Run executes the seed application.
func (s *Seed) Run(ctx context.Context) error {
	// Set output token tpl
	s.SetOutputTokenTpl()

	// Output config
	outputConfig, err := s.OutputConfig(ctx)
	if err != nil {
		return fmt.Errorf("output config: %w", err)
	}

	// Create admin static token
	adminStaticToken, err := s.CreateAdminStaticToken(ctx)
	if err != nil {
		return fmt.Errorf("create admin static token: %w", err)
	}

	// Create admin role
	if err := s.CreateAdminRole(ctx); err != nil {
		return fmt.Errorf("create admin role: %w", err)
	}

	// Create admin http rule
	if err := s.CreateAdminHTTPRule(ctx); err != nil {
		return fmt.Errorf("create admin http rule: %w", err)
	}

	// Output
	if err := s.Output(
		outputConfig,
		&OutputData{
			AdminStaticToken: adminStaticToken,
		},
	); err != nil {
		return fmt.Errorf("output exec: %w", err)
	}

	s.Logger.InfoContext(ctx, "seed completed")

	return nil
}

// CreateAdminStaticToken creates or retrieves the admin static token.
func (s *Seed) CreateAdminStaticToken(ctx context.Context) (string, error) {
	adminStaticToken := ""

	// Filter admin static token
	if tokens, err := s.Service.Tokens.FilterStaticAccessTokens(
		ctx,
		tokenssp.FilterStaticAccessTokenData{
			ID: &[]string{s.Config.Admin.Token.ID},
		},
	); err != nil {
		return "", fmt.Errorf("failed to filter admin static token: %w", err)
	} else if len(tokens) == 0 {
		// Create admin static token
		token, err := s.Service.Tokens.CreateStaticAccessToken(
			ctx,
			tokenssp.CreateStaticAccessTokenData{
				ID:          s.Config.Admin.Token.ID,
				Roles:       []string{s.Config.Admin.Role.ID},
				Description: s.Config.Admin.Token.Description,
			},
		)
		if err != nil {
			return "", fmt.Errorf("failed to create admin static token: %w", err)
		}

		s.Logger.InfoContext(
			ctx,
			"admin static token created",
			slog.String("id", s.Config.Admin.Token.ID),
			slog.String("role", s.Config.Admin.Role.ID),
			slog.String("description", s.Config.Admin.Token.Description),
		)

		adminStaticToken = token
	}

	return adminStaticToken, nil
}

// CreateAdminRole creates the admin role if it does not exist.
func (s *Seed) CreateAdminRole(ctx context.Context) error {
	// Filter admin role
	if roles, err := s.Service.Roles.FilterRoles(
		ctx,
		rolessp.FilterRolesData{
			ID:          &[]string{s.Config.Admin.Role.ID},
			Name:        nil,
			SystemFlag:  nil,
			ServiceFlag: nil,
		},
	); err != nil {
		return fmt.Errorf("failed to filter admin role: %w", err)
	} else if len(roles) == 0 {
		// Create admin role
		if _, err := s.Service.Roles.CreateRole(
			ctx,
			rolessp.CreateRoleData{
				ID:          s.Config.Admin.Role.ID,
				Name:        s.Config.Admin.Role.Name,
				Description: s.Config.Admin.Role.Description,
				SystemFlag:  true,
				ServiceFlag: false,
			},
		); err != nil {
			return fmt.Errorf("failed to create admin role: %w", err)
		}

		s.Logger.InfoContext(
			ctx,
			"admin role created",
			slog.String("id", s.Config.Admin.Role.ID),
			slog.String("name", s.Config.Admin.Role.Name),
			slog.String("description", s.Config.Admin.Role.Description),
			slog.Bool("system_flag", true),
			slog.Bool("service_flag", false),
		)
	}

	return nil
}

// CreateAdminHTTPRule creates the admin HTTP rule if it does not exist.
func (s *Seed) CreateAdminHTTPRule(ctx context.Context) error {
	// Filter admin http rule
	if rules, err := s.Service.Rules.FilterHTTPRules(
		ctx,
		rulessp.FilterHTTPRulesData{
			ID:      nil,
			RoleID:  &[]string{s.Config.Admin.Role.ID},
			Path:    &[]string{s.Config.Admin.Rule.Path},
			Methods: &s.Config.Admin.Rule.Methods,
			Mfa:     func(b bool) *bool { return &b }(false),
		},
	); err != nil {
		return fmt.Errorf("failed to filter admin http rule: %w", err)
	} else if len(rules) == 0 {
		// Create admin http rule
		if _, err := s.Service.Rules.CreateHTTPRule(
			ctx,
			rulessp.CreateHTTPRuleData{
				RoleID:  s.Config.Admin.Role.ID,
				Path:    s.Config.Admin.Rule.Path,
				Methods: s.Config.Admin.Rule.Methods,
				Mfa:     false,
			},
		); err != nil {
			return fmt.Errorf("failed to create admin http rule: %w", err)
		}

		s.Logger.InfoContext(
			ctx,
			"admin http rule created",
			slog.String("role", s.Config.Admin.Role.ID),
			slog.String("path", s.Config.Admin.Rule.Path),
			slog.String("methods", strings.Join(s.Config.Admin.Rule.Methods, ",")),
			slog.Bool("mfa", false),
		)
	}

	return nil
}

// OutputConfig prepares writers for outputting the admin token.
func (s *Seed) OutputConfig(ctx context.Context) (*OutputConfig, error) {
	if s.Options.Output.JSON {
		return s.outputConfigJSON(ctx)
	}

	switch {
	case s.Options.Output.AdminToken == "" || s.Options.Output.AdminToken == "stdout":
		return s.outputConfigStdout(ctx), nil
	case strings.HasPrefix(s.Options.Output.AdminToken, "file:"):
		path := strings.TrimPrefix(
			s.Options.Output.AdminToken,
			"file:",
		)

		return s.outputConfigFile(ctx, path)

	default:
		return nil, fmt.Errorf(
			"%w: %v",
			ErrInvalidAdminTokenOutput,
			s.Options.Output.AdminToken,
		)
	}
}

// Output prints or writes the generated admin token based on output options.
func (s *Seed) Output(config *OutputConfig, data *OutputData) error {
	// Output json
	if s.Options.Output.JSON {
		if err := s.outputJSON(data); err != nil {
			return fmt.Errorf("output json: %w", err)
		}

		return nil
	}

	// Output custom
	if err := s.outputCustom(config, data); err != nil {
		return fmt.Errorf("output custom: %w", err)
	}

	return nil
}

// OutputQuietWrite prints a value quietly to the given writer.
func (s *Seed) OutputQuietWrite(w io.Writer, str string) error {
	fn := fmt.Fprintln
	if w != s.Options.Output.Stdout {
		fn = fmt.Fprint
	}

	if _, err := fn(w, str); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}

// SetOutputTokenTpl initializes the template used to render
// a one-time display of a sensitive token with a security warning.
func (s *Seed) SetOutputTokenTpl() {
	s.outputTokenTpl = template.Must(template.New("outputTokenTpl").Parse(
		`{{.Sep}}
SECURITY WARNING!

The following output contains highly sensitive information intended for the operator only.` +
			`This data is displayed ONCE and must be stored securely. Do NOT share it, commit it ` +
			`to version control, or expose it in logs, CI systems, or monitoring tools.

{{.Token}}
{{.Sep}}`,
	))
}

func (s *Seed) outputConfigJSON(
	ctx context.Context,
) (*OutputConfig, error) {
	s.Logger.InfoContext(
		ctx,
		"outputting in JSON format to stdout",
	)

	return &OutputConfig{
		AdminTokenWriter: nil,
	}, nil
}

func (s *Seed) outputConfigStdout(ctx context.Context) *OutputConfig {
	s.Logger.InfoContext(
		ctx,
		"admin static token displayed in stdout",
		slog.Bool("quiet", s.Options.Output.Quiet),
	)

	return &OutputConfig{
		AdminTokenWriter: s.Options.Output.Stdout,
	}
}

func (s *Seed) outputConfigFile(ctx context.Context, path string) (*OutputConfig, error) {
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

	if err := shutdown.AddHandler(
		func(ctx context.Context, _ int) error {
			errCh := make(chan error, 1)

			go func() {
				errCh <- f.Close()
			}()

			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-errCh:
				return err
			}
		},
	); err != nil {
		return nil, fmt.Errorf("add shutdown handler: %w", err)
	}

	s.Logger.InfoContext(
		ctx,
		"admin static token written to file",
		slog.String("path", path),
		slog.Bool("quiet", s.Options.Output.Quiet),
	)

	return &OutputConfig{AdminTokenWriter: f}, nil
}

func (s *Seed) outputJSON(data *OutputData) error {
	enc := json.NewEncoder(s.Options.Output.Stdout)
	if s.Options.Output.JSONPretty {
		enc.SetIndent("", "  ")
	}

	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("marshal output data: %w", err)
	}

	return nil
}

func (s *Seed) outputCustom(config *OutputConfig, data *OutputData) error {
	if s.Options.Output.Quiet {
		var token string
		if data.AdminStaticToken != "" {
			token = data.AdminStaticToken
		} else {
			token = noStaticAdminToken
		}

		if err := s.OutputQuietWrite(
			config.AdminTokenWriter,
			token,
		); err != nil {
			return fmt.Errorf("output quiet print: admin static token: %w", err)
		}

		return nil
	}

	var token string

	if data.AdminStaticToken != "" {
		token = data.AdminStaticToken
	} else {
		token = noStaticAdminToken
	}

	if err := s.outputTokenTpl.Execute(
		config.AdminTokenWriter,
		outputTokenData{
			Sep:   strings.Repeat("-", termWidth()),
			Token: token,
		},
	); err != nil {
		return fmt.Errorf("template execute: %w", err)
	}

	return nil
}

// termWidth returns the terminal width or defaults to 80 if unavailable.
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
