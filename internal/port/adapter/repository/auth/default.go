// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// JWTTokenAudience is the JWT audience claim used for tokens in the auth service.
	JWTTokenAudience = "auth-service"

	// JWTSignKeyMinLen is the minimum key length in bytes for signing JWT tokens
	// (HS512, 512 bits).
	JWTSignKeyMinLen = 64

	// JWTHashKeyMinLen is the minimum key length in bytes for hashing JWT tokens
	// (HMAC-SHA-256, 256 bits).
	JWTHashKeyMinLen = 32

	// JWTLeeway is the duration used to account for clock skew during token validation.
	JWTLeeway = 5 * time.Second

	// StaticTokenUser is the user ID used for system/static tokens.
	StaticTokenUser = 0

	// StaticTokenDevice is the device identifier used for system/static tokens.
	StaticTokenDevice = "system"

	// AuthKeyMinLen is the minimum key length in bytes for AES-256-GCM encryption
	// of authentication requests and responses.
	AuthKeyMinLen = 32
)

// JWTSigningMethod is the JWT signing method used for token creation and verification.
func JWTSigningMethod() jwt.SigningMethod {
	return jwt.SigningMethodHS512
}
