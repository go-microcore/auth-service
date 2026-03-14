// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package auth

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.microcore.dev/auth-service/internal/adapter/repository/auth/model"
	"go.microcore.dev/auth-service/internal/port/adapter/repository/auth"
	"go.microcore.dev/auth-service/internal/shared/postgres"
	"go.microcore.dev/auth-service/internal/shared/redis"
)

// ErrCiphertextTooShort is returned when the ciphertext length is smaller than the
// expected nonce size.
var ErrCiphertextTooShort = errors.New("ciphertext too short")

type (
	// AdapterConfig provides auth repository adapter configuration.
	AdapterConfig struct {
		Config          *Config
		Logger          *slog.Logger
		Redis           *redis.Redis
		Postgres        *postgres.Postgres
		LocalTokenCache *LocalTokenCache
	}

	adapter struct {
		*AdapterConfig
	}

	// jwtTokenClaims defines custom JWT claims for a user, including device, roles, and MFA status.
	jwtTokenClaims struct {
		jwt.RegisteredClaims          // Standard JWT claims (iss, exp, iat, etc.)
		Device               string   `json:"device"` // Device identifier
		Roles                []string `json:"roles"`  // User roles
		Mfa                  bool     `json:"mfa"`    // Indicates if MFA was used
	}
)

// NewAdapter creates a new instance of the adapter.
func NewAdapter(config *AdapterConfig) auth.Adapter {
	return &adapter{config}
}

// DecryptAuthRequest decrypt auth request.
func (a *adapter) DecryptAuthRequest(
	ctx context.Context,
	data []byte,
) ([]byte, error) {
	block, err := aes.NewCipher(a.Config.Auth.Key)
	if err != nil {
		return nil, fmt.Errorf("aes new cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cipher new gcm: %w", err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return nil, ErrCiphertextTooShort
	}

	nonce, data := data[:nonceSize], data[nonceSize:]

	res, err := aesGCM.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, fmt.Errorf("aes gcm open: %w", err)
	}

	return res, nil
}

// EncryptAuthResponse encrypt auth response.
func (a *adapter) EncryptAuthResponse(
	ctx context.Context,
	data []byte,
) ([]byte, error) {
	block, err := aes.NewCipher(a.Config.Auth.Key)
	if err != nil {
		return nil, fmt.Errorf("aes new cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cipher new gcm: %w", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("io read full: %w", err)
	}

	ciphertext := aesGCM.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

// ParseAccessToken parses a JWT access token.
func (a *adapter) ParseAccessToken(
	ctx context.Context,
	token string,
) (*auth.ParseTokenResult, error) {
	return a.parseToken(ctx, token, a.Config.JWT.AccessKey)
}

// ParseRefreshToken parses a JWT refresh token.
func (a *adapter) ParseRefreshToken(
	ctx context.Context,
	token string,
) (*auth.ParseTokenResult, error) {
	return a.parseToken(ctx, token, a.Config.JWT.RefreshKey)
}

// NewTokens generates new access and refresh tokens.
func (a *adapter) NewTokens(
	ctx context.Context,
	data auth.NewTokenData,
) (*auth.NewTokensResult, error) {
	// Generate JWT Access token
	access, err := a.newAccessToken(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("access: %w", err)
	}

	// Generate JWT Refresh token
	refresh, err := a.newRefreshToken(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("refresh: %w", err)
	}

	return &auth.NewTokensResult{
		Access:  access,
		Refresh: refresh,
	}, nil
}

// CreateStaticAccessToken creates a static access token.
func (a *adapter) CreateStaticAccessToken(
	ctx context.Context,
	data auth.CreateStaticAccessTokenData,
) (string, error) {
	// Generate JWT static access token
	access, err := a.newToken(
		ctx,
		auth.NewTokenData{
			User:   auth.StaticTokenUser,
			Roles:  data.Roles,
			Mfa:    false,
			Device: auth.StaticTokenDevice,
		},
		a.Config.JWT.AccessKey,
		0,
	)
	if err != nil {
		return "", fmt.Errorf("new token: %w", err)
	}

	token, err := a.jwtHash(access)
	if err != nil {
		return "", fmt.Errorf("jwt hash: %w", err)
	}

	// Create token model
	obj := model.AuthStaticToken{
		ID:          data.ID,
		Token:       token,
		UserID:      auth.StaticTokenUser,
		Device:      auth.StaticTokenDevice,
		Roles:       data.Roles,
		Description: data.Description,
		Created:     time.Unix(0, time.Now().UnixNano()),
	}

	// Save token to database
	if err := a.Postgres.Manager.Client().
		WithContext(ctx).
		Create(&obj).
		Error; err != nil {
		return "", fmt.Errorf("pg: %w", err)
	}

	return access, nil
}

// NewSession creates a new user session.
func (a *adapter) NewSession(
	ctx context.Context,
	user uint,
	device string,
	session *auth.Session,
) error {
	key := fmt.Sprintf("users:sessions:%d:%s", user, device)
	devicesKey := fmt.Sprintf("users:devices:%d", user)

	// Set session
	if err := a.Redis.Manager.Client().HSet(ctx, key, map[string]any{
		"jti":             session.Jti,
		"issued_at":       session.IssuedAt,
		"location":        session.Location,
		"ip":              session.IP,
		"user_agent":      session.UserAgent,
		"os_full_name":    session.OsFullName,
		"os_name":         session.OsName,
		"os_version":      session.OsVersion,
		"platform":        session.Platform,
		"model":           session.Model,
		"browser_name":    session.BrowserName,
		"browser_version": session.BrowserVersion,
		"engine_name":     session.EngineName,
		"engine_version":  session.EngineVersion,
	}).Err(); err != nil {
		return fmt.Errorf("set: %w", err)
	}

	// Expire session
	if err := a.Redis.Manager.Client().
		Expire(ctx, key, a.Config.JWT.RefreshTTL).
		Err(); err != nil {
		return fmt.Errorf("expire: %w", err)
	}

	// Add device
	if err := a.Redis.Manager.Client().
		SAdd(ctx, devicesKey, device).
		Err(); err != nil {
		return fmt.Errorf("add device: %w", err)
	}

	return nil
}

// UpdateSession updates a session's JTI and TTL.
func (a *adapter) UpdateSession(
	ctx context.Context,
	user uint,
	device, jti string,
) error {
	key := fmt.Sprintf("users:sessions:%d:%s", user, device)

	// Update jti
	if err := a.Redis.Manager.Client().
		HSet(ctx, key, "jti", jti).
		Err(); err != nil {
		return fmt.Errorf("set jti: %w", err)
	}

	// Update ttl
	if err := a.Redis.Manager.Client().
		Expire(ctx, key, a.Config.JWT.RefreshTTL).
		Err(); err != nil {
		return fmt.Errorf("expire: %w", err)
	}

	return nil
}

// DeleteSession deletes a user session and device.
func (a *adapter) DeleteSession(
	ctx context.Context,
	user uint,
	device string,
) error {
	sessionKey := fmt.Sprintf("users:sessions:%d:%s", user, device)
	devicesKey := fmt.Sprintf("users:devices:%d", user)

	// Delete session
	if err := a.Redis.Manager.Client().
		Del(ctx, sessionKey).
		Err(); err != nil {
		return fmt.Errorf("delele: %w", err)
	}

	// Delete device
	if err := a.Redis.Manager.Client().
		SRem(ctx, devicesKey, device).
		Err(); err != nil {
		return fmt.Errorf("delete device: %w", err)
	}

	return nil
}

// GetSession retrieves a user session by device.
func (a *adapter) GetSession(
	ctx context.Context,
	user uint,
	device string,
) (*auth.Session, error) {
	key := fmt.Sprintf("users:sessions:%d:%s", user, device)

	data, err := a.Redis.Manager.
		Client().
		HGetAll(ctx, key).
		Result()
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	if len(data) == 0 {
		return nil, auth.ErrSessionNotFound
	}

	return &auth.Session{
		Jti:            data["jti"],
		IssuedAt:       data["issued_at"],
		Location:       data["location"],
		IP:             data["ip"],
		UserAgent:      data["user_agent"],
		OsFullName:     data["os_full_name"],
		OsName:         data["os_name"],
		OsVersion:      data["os_version"],
		Platform:       data["platform"],
		Model:          data["model"],
		BrowserName:    data["browser_name"],
		BrowserVersion: data["browser_version"],
		EngineName:     data["engine_name"],
		EngineVersion:  data["engine_version"],
	}, nil
}

// GetActiveDevices returns active devices for a user.
func (a *adapter) GetActiveDevices(
	ctx context.Context,
	user uint,
) ([]auth.Device, error) {
	devicesKey := fmt.Sprintf("users:devices:%d", user)

	devices, err := a.Redis.Manager.
		Client().
		SMembers(ctx, devicesKey).
		Result()
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	res := make([]auth.Device, len(devices))
	for i := range devices {
		session, err := a.GetSession(ctx, user, devices[i])
		if err != nil && !errors.Is(err, auth.ErrSessionNotFound) {
			return nil, fmt.Errorf("get session: %w", err)
		}

		if session == nil {
			// Remove expired device from Set
			_ = a.Redis.Manager.Client().SRem(ctx, devicesKey, devices[i])
		} else {
			res[i] = auth.Device{
				ID:      devices[i],
				Session: *session,
			}
		}
	}

	return res, nil
}

// FilterStaticAccessTokens returns static access tokens filtered by ID.
func (a *adapter) FilterStaticAccessTokens(
	ctx context.Context,
	data auth.FilterStaticAccessTokenData,
) ([]auth.StaticAccessTokenResult, error) {
	// Create model
	tokens := []model.AuthStaticToken{}

	// Create query with context
	query := a.Postgres.Manager.Client().WithContext(ctx)

	// Filter by id
	if data.ID != nil {
		query = query.Where("id IN ?", *data.ID)
	}

	// Get tokens from database
	if err := query.Find(&tokens).Error; err != nil {
		return nil, fmt.Errorf("pg: %w", err)
	}

	// Make res
	res := make([]auth.StaticAccessTokenResult, len(tokens))
	for i := range tokens {
		res[i] = auth.StaticAccessTokenResult{
			ID:          tokens[i].ID,
			Token:       tokens[i].Token,
			UserID:      tokens[i].UserID,
			Device:      tokens[i].Device,
			Roles:       []string(tokens[i].Roles),
			Description: tokens[i].Description,
			Created:     tokens[i].Created,
		}
	}

	return res, nil
}

// DeleteStaticAccessToken deletes a static access token by ID.
func (a *adapter) DeleteStaticAccessToken(
	ctx context.Context,
	id string,
) error {
	// Delete token from database
	result := a.Postgres.Manager.Client().
		WithContext(ctx).
		Delete(&model.AuthStaticToken{}, "id = ?", id)

	// Check errors
	if result.Error != nil {
		return fmt.Errorf("pg: %w", result.Error)
	}

	// If token not found
	if result.RowsAffected == 0 {
		return auth.ErrStaticTokenNotFound
	}

	return nil
}

// newToken generates a JWT token with given claims.
func (a *adapter) newToken(
	_ context.Context,
	data auth.NewTokenData,
	key []byte,
	expire time.Duration,
) (string, error) {
	jti := uuid.NewString()
	now := time.Now()

	// Create claims
	claims := &jwtTokenClaims{
		Device: data.Device,
		Roles:  data.Roles,
		Mfa:    data.Mfa,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:       jti,
			Subject:  strconv.FormatUint(uint64(data.User), 10),
			IssuedAt: jwt.NewNumericDate(now),
			Issuer:   a.Config.JWT.Issuer,
			Audience: []string{auth.JWTTokenAudience},
		},
	}

	if expire > 0 {
		claims.ExpiresAt = jwt.NewNumericDate(now.Add(expire))
	}

	// Create token
	token := jwt.NewWithClaims(auth.JWTSigningMethod(), claims)

	// Create signed token string
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return tokenString, nil
}

// newAccessToken generates a new JWT access token.
func (a *adapter) newAccessToken(
	ctx context.Context,
	data auth.NewTokenData,
) (string, error) {
	return a.newToken(ctx, data, a.Config.JWT.AccessKey, a.Config.JWT.AccessTTL)
}

// newRefreshToken generates a new JWT refresh token.
func (a *adapter) newRefreshToken(
	ctx context.Context,
	data auth.NewTokenData,
) (string, error) {
	return a.newToken(ctx, data, a.Config.JWT.RefreshKey, a.Config.JWT.RefreshTTL)
}

// parseToken parses a JWT or static token and validates it.
//
//nolint:gocognit,cyclop,funlen // need refactoring
func (a *adapter) parseToken(
	ctx context.Context,
	token string,
	key []byte,
) (*auth.ParseTokenResult, error) {
	tokenHash, err := a.jwtHash(token)
	if err != nil {
		return nil, fmt.Errorf("jwt hash: %w", err)
	}

	tokenKey := "auth:tokens:" + tokenHash
	tokenTTL := a.Config.Cache.StaticTokenTTL // default for static tokens

	// Get local cache
	if res, ok := a.LocalTokenCache.Get(tokenKey); ok {
		return res, nil
	}

	// Get token from cache
	if data, err := a.Redis.Manager.
		Client().
		Get(ctx, tokenKey).
		Result(); err == nil {
		var result auth.ParseTokenResult
		if err := json.Unmarshal([]byte(data), &result); err == nil {
			if result.Expires != nil {
				tokenTTL = time.Until(time.Unix(*result.Expires, 0))
			}

			a.LocalTokenCache.Set(tokenKey, &result, tokenTTL)

			return &result, nil
		}
	}

	// Creating a parser with options
	parser := jwt.NewParser(
		// Manual security control that disallows the use of any algorithms other than HMAC.
		// You explicitly state that only HMAC (e.g., HS256, HS512) is accepted and everything
		// else is rejected. This check explicitly forbids any unsupported algorithms.
		jwt.WithValidMethods([]string{auth.JWTSigningMethod().Alg()}),
		// In many distributed systems, clocks may differ slightly between services (for example,
		// the client and server might have a small time difference). To avoid issues caused by
		// this, we add a leeway (a time buffer).
		jwt.WithLeeway(auth.JWTLeeway),
		jwt.WithIssuer(a.Config.JWT.Issuer),
		jwt.WithAudience(auth.JWTTokenAudience),
		jwt.WithIssuedAt(),
	)

	// Create token claims
	claims := &jwtTokenClaims{}

	// Parsing the token with custom claims
	jwtToken, err := parser.ParseWithClaims(
		token,
		claims,
		func(_ *jwt.Token) (any, error) {
			return key, nil
		},
	)
	if err != nil || !jwtToken.Valid {
		return nil, auth.ErrInvalidToken
	}

	// Check expires
	var expires *int64

	if claims.ExpiresAt == nil {
		// Static token
		var count int64

		err := a.Postgres.Manager.Client().WithContext(ctx).
			Model(&model.AuthStaticToken{}).
			Where("token = ?", tokenHash).
			Count(&count).Error
		if err != nil {
			return nil, fmt.Errorf("pg: %w", err)
		}

		if count == 0 {
			return nil, auth.ErrInvalidToken
		}
	} else {
		// Regular token
		e := claims.ExpiresAt.Unix()

		expires = &e
	}

	// Token TTL for regular tokens
	if expires != nil {
		tokenTTL = time.Until(time.Unix(*expires, 0))
		if tokenTTL <= 0 {
			return nil, auth.ErrInvalidToken
		}
	}

	// Parsing the Subject into a uint
	user, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("subject parse: %w", err)
	}

	result := &auth.ParseTokenResult{
		ID:       claims.ID,
		Device:   claims.Device,
		User:     uint(user),
		Roles:    claims.Roles,
		Mfa:      claims.Mfa,
		Expires:  expires,
		Issued:   claims.IssuedAt.Unix(),
		Issuer:   claims.Issuer,
		Audience: claims.Audience,
	}

	a.LocalTokenCache.Set(tokenKey, result, tokenTTL)

	// Save token to cache
	b, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("json marshal result: %w", err)
	}

	if err := a.Redis.Manager.Client().
		Set(ctx, tokenKey, b, tokenTTL).
		Err(); err != nil {
		return nil, fmt.Errorf("set redis cache: %w", err)
	}

	return result, nil
}

// jwtHash computes HMAC-SHA256 hash of a token.
func (a *adapter) jwtHash(token string) (string, error) {
	mac := hmac.New(sha256.New, a.Config.JWT.HashKey)
	if _, err := mac.Write([]byte(token)); err != nil {
		return "", fmt.Errorf("mac write: %w", err)
	}

	return hex.EncodeToString(mac.Sum(nil)), nil
}
