// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package tokens

import (
	"time"
)

type (
	// Data

	// AuthData holds information needed to authenticate a user and create tokens.
	AuthData struct {
		User               uint
		Roles              []string
		Mfa                bool
		Device             string
		MetaLocation       string
		MetaIP             string
		MetaUserAgent      string
		MetaOsFullName     string
		MetaOsName         string
		MetaOsVersion      string
		MetaPlatform       string
		MetaModel          string
		MetaBrowserName    string
		MetaBrowserVersion string
		MetaEngineName     string
		MetaEngineVersion  string
		TTL                time.Time
	}

	// Auth2faData holds information for two-factor authentication.
	Auth2faData struct {
		User   uint
		Roles  []string
		Device string
		TTL    time.Time
	}

	// TokenRenewData holds data needed to renew a refresh token.
	TokenRenewData struct {
		RefreshToken string
	}

	// TokenValidateData holds data needed to validate an access token.
	TokenValidateData struct {
		AccessToken string
	}

	// CreateStaticAccessTokenData holds data to create a static access token.
	CreateStaticAccessTokenData struct {
		ID          string
		Roles       []string
		Description string
	}

	// FilterStaticAccessTokenData holds criteria to filter static access tokens.
	FilterStaticAccessTokenData struct {
		ID *[]string
	}

	// Results

	// AuthResult contains tokens and session info after authentication.
	AuthResult struct {
		Access    string
		Refresh   string
		Mfa       bool
		NewDevice bool
	}

	// Auth2faResult contains tokens after completing two-factor authentication.
	Auth2faResult struct {
		Access  string
		Refresh string
	}

	// TokenRenewResult contains renewed access and refresh tokens.
	TokenRenewResult struct {
		Access  string
		Refresh string
		Mfa     bool
	}

	// TokenValidateResult contains decoded information from a validated access token.
	TokenValidateResult struct {
		ID       string
		Device   string
		User     uint
		Roles    []string
		Mfa      bool
		Expires  *int64
		Issued   int64
		Issuer   string
		Audience []string
	}

	// StaticAccessTokenResult represents a stored static access token with metadata.
	StaticAccessTokenResult struct {
		ID          string
		Token       string
		UserID      uint
		Device      string
		Roles       []string
		Description string
		Created     time.Time
	}
)
