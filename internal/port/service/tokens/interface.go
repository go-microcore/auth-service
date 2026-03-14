// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package tokens

import (
	"context"
)

type (
	// Service defines the interface for the tokens service.
	Service interface {
		DecryptAuthRequest(
			ctx context.Context,
			data []byte,
			target any,
		) error

		EncryptAuthResponse(
			ctx context.Context,
			data []byte,
		) ([]byte, error)

		Auth(
			ctx context.Context,
			data *AuthData,
		) (*AuthResult, error)

		Auth2fa(
			ctx context.Context,
			data Auth2faData,
		) (*Auth2faResult, error)

		TokenRenew(
			ctx context.Context,
			data TokenRenewData,
		) (*TokenRenewResult, error)

		TokenValidate(
			ctx context.Context,
			data TokenValidateData,
		) (*TokenValidateResult, error)

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
