// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package auth

import (
	"time"
)

type (
	// Data

	// CreateStaticAccessTokenData represents the input data for creating a static access token.
	CreateStaticAccessTokenData struct {
		ID          string
		Roles       []string
		Description string
	}

	// FilterStaticAccessTokenData represents the filter parameters for querying static access tokens.
	FilterStaticAccessTokenData struct {
		ID *[]string
	}

	// NewTokenData represents the data required to generate new access and refresh tokens.
	NewTokenData struct {
		User   uint
		Roles  []string
		Mfa    bool
		Device string
	}

	// Results

	// ParseTokenResult contains information extracted from a parsed JWT token.
	ParseTokenResult struct {
		ID       string   `json:"id"`
		Device   string   `json:"device"`
		User     uint     `json:"user"`
		Roles    []string `json:"roles"`
		Mfa      bool     `json:"mfa"`
		Expires  *int64   `json:"expires"`
		Issued   int64    `json:"issued"`
		Issuer   string   `json:"issuer"`
		Audience []string `json:"audience"`
	}

	// NewTokensResult contains newly generated access and refresh tokens.
	NewTokensResult struct {
		Access  string
		Refresh string
	}

	// Session represents a user session on a device, including metadata about
	// the device and environment.
	Session struct {
		Jti            string
		IssuedAt       string
		Location       string
		IP             string
		UserAgent      string
		OsFullName     string
		OsName         string
		OsVersion      string
		Platform       string
		Model          string
		BrowserName    string
		BrowserVersion string
		EngineName     string
		EngineVersion  string
	}

	// Device represents a user device and its associated session.
	Device struct {
		ID      string
		Session Session
	}

	// StaticAccessTokenResult represents a stored static access token and its metadata.
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
