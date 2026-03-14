// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package auth

import (
	"context"
)

type (
	// Adapter defines the interface for the auth repository adapter.
	Adapter interface {
		DecryptAuthRequest(
			ctx context.Context,
			data []byte,
		) ([]byte, error)

		EncryptAuthResponse(
			ctx context.Context,
			data []byte,
		) ([]byte, error)

		ParseAccessToken(
			ctx context.Context,
			token string,
		) (*ParseTokenResult, error)

		ParseRefreshToken(
			ctx context.Context,
			token string,
		) (*ParseTokenResult, error)

		NewTokens(
			ctx context.Context,
			data NewTokenData,
		) (*NewTokensResult, error)

		NewSession(
			ctx context.Context,
			user uint,
			device string,
			session *Session,
		) error

		UpdateSession(
			ctx context.Context,
			user uint,
			device string,
			jti string,
		) error

		DeleteSession(
			ctx context.Context,
			user uint,
			device string,
		) error

		GetSession(
			ctx context.Context,
			user uint,
			device string,
		) (*Session, error)

		GetActiveDevices(
			ctx context.Context,
			user uint,
		) ([]Device, error)

		// Static access tokens
		CreateStaticAccessToken(
			ctx context.Context,
			data CreateStaticAccessTokenData,
		) (string, error)

		FilterStaticAccessTokens(
			ctx context.Context,
			data FilterStaticAccessTokenData,
		) ([]StaticAccessTokenResult, error)

		DeleteStaticAccessToken(
			ctx context.Context,
			id string,
		) error
	}
)
