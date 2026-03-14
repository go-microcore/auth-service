// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package http

import (
	"time"
)

type (
	// Requests

	// AuthRequest represents a request to get an access and refresh tokens.
	AuthRequest struct {
		User               uint      `json:"user"`
		Roles              []string  `json:"roles"`
		Mfa                bool      `json:"mfa"`
		Device             string    `json:"device"`
		MetaLocation       string    `json:"metaLocation"`
		MetaIP             string    `json:"metaIp"`
		MetaUserAgent      string    `json:"metaUserAgent"`
		MetaOsFullName     string    `json:"metaOsFullName"`
		MetaOsName         string    `json:"metaOsName"`
		MetaOsVersion      string    `json:"metaOsVersion"`
		MetaPlatform       string    `json:"metaPlatform"`
		MetaModel          string    `json:"metaModel"`
		MetaBrowserName    string    `json:"metaBrowserName"`
		MetaBrowserVersion string    `json:"metaBrowserVersion"`
		MetaEngineName     string    `json:"metaEngineName"`
		MetaEngineVersion  string    `json:"metaEngineVersion"`
		TTL                time.Time `json:"ttl"`
	}

	// Auth2FARequest represents a request to get an access and refresh tokens
	// after successful 2FA.
	Auth2FARequest struct {
		User   uint      `json:"user"`
		Roles  []string  `json:"roles"`
		Device string    `json:"device"`
		TTL    time.Time `json:"ttl"`
	}

	// TokenRenewRequest represents a request to renew an access token using a refresh token.
	TokenRenewRequest struct {
		RefreshToken string `json:"refreshToken"`
	}

	// TokenAuthorizeHTTPRequest represents a request to authorize an HTTP request by path and method.
	TokenAuthorizeHTTPRequest struct {
		Path   string `json:"path"`
		Method string `json:"method"`
	}

	// CreateStaticAccessTokenRequest represents a request to create a static access token.
	CreateStaticAccessTokenRequest struct {
		ID          string   `json:"id"`
		Roles       []string `json:"roles"`
		Description string   `json:"description"`
	}

	// FilterStaticAccessTokenRequest represents a request to filter static tokens by ID.
	FilterStaticAccessTokenRequest struct {
		ID *[]string `json:"id"`
	}

	// Responses

	// AuthResponse represents the response returned after a successful authentication.
	AuthResponse struct {
		Access    string `json:"access"`
		Refresh   string `json:"refresh"`
		Mfa       bool   `json:"mfa"`
		NewDevice bool   `json:"newDevice"`
	}

	// Auth2FAResponse represents a 2FA authentication response.
	Auth2FAResponse struct {
		Access  string `json:"access"`
		Refresh string `json:"refresh"`
	}

	// TokenRenewResponse represents the response when renewing tokens.
	TokenRenewResponse struct {
		Access  string `json:"accessToken"`
		Refresh string `json:"refreshToken"`
		Mfa     bool   `json:"mfaRequired"`
	}

	// TokenValidateResponse represents the response for token validation.
	TokenValidateResponse struct {
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

	// TokenAuthorizeHTTPResponse represents the response for HTTP token authorization.
	TokenAuthorizeHTTPResponse struct {
		Token TokenAuthorizeHTTPDataResponse `json:"token"`
		Auth  TokenAuthorizeHTTPAuthResponse `json:"auth"`
	}

	// TokenAuthorizeHTTPDataResponse contains the token data in an HTTP authorization response.
	TokenAuthorizeHTTPDataResponse struct {
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

	// TokenAuthorizeHTTPAuthResponse contains the auth info in an HTTP authorization response.
	TokenAuthorizeHTTPAuthResponse struct {
		Mfa bool `json:"mfa"`
	}

	// CreateStaticAccessTokenResponse represents the response after creating a static access token.
	CreateStaticAccessTokenResponse struct {
		Token string `json:"token"`
	}

	// FilterStaticAccessTokenResponse represents a filtered static token and its metadata.
	FilterStaticAccessTokenResponse struct {
		ID          string    `json:"id"`
		Token       string    `json:"token"`
		UserID      uint      `json:"userId"`
		Device      string    `json:"device"`
		Roles       []string  `json:"roles"`
		Description string    `json:"description"`
		Created     time.Time `json:"created"`
	}
)

// Validate checks that the refresh token is not empty.
func (r *TokenRenewRequest) Validate() error {
	if r.RefreshToken == "" {
		return ErrInvalidToken
	}

	return nil
}

// Validate checks that both path and method are set.
func (r *TokenAuthorizeHTTPRequest) Validate() error {
	if r.Path == "" {
		return ErrInvalidPath
	}

	if r.Method == "" {
		return ErrInvalidMethod
	}

	return nil
}

// Validate checks all fields of the request.
func (r *CreateStaticAccessTokenRequest) Validate() error {
	if r.ID == "" {
		return ErrInvalidID
	}

	if len(r.Roles) == 0 {
		return ErrInvalidRoles
	}

	if r.Description == "" {
		return ErrInvalidDescription
	}

	return nil
}
