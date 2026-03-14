// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package auth

import (
	"fmt"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
	port "go.microcore.dev/auth-service/internal/port/adapter/repository/auth"
)

type (
	// LocalTokenCache is a thread-safe LRU cache for parsed tokens.
	LocalTokenCache struct {
		cache *lru.Cache
		mu    sync.Mutex
	}

	// tokenCacheItem stores a cached token and its expiration time.
	tokenCacheItem struct {
		result  *port.ParseTokenResult
		expires time.Time
	}
)

// NewLocalTokenCache creates a new LocalTokenCache with the given size.
func NewLocalTokenCache(size int) (*LocalTokenCache, error) {
	c, err := lru.New(size)
	if err != nil {
		return nil, fmt.Errorf("new loacl token cache: %w", err)
	}

	return &LocalTokenCache{
		cache: c,
		mu:    sync.Mutex{},
	}, nil
}

// Get retrieves a cached token by key if it exists and is not expired.
func (c *LocalTokenCache) Get(key string) (*port.ParseTokenResult, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if val, ok := c.cache.Get(key); ok {
		item, ok := val.(tokenCacheItem)
		if !ok {
			c.cache.Remove(key)
			return nil, false
		}

		if time.Now().Before(item.expires) {
			return item.result, true
		}

		// expired, remove
		c.cache.Remove(key)
	}

	return nil, false
}

// Set stores a token in the cache with a TTL.
func (c *LocalTokenCache) Set(key string, result *port.ParseTokenResult, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache.Add(key, tokenCacheItem{
		result:  result,
		expires: time.Now().Add(ttl),
	})
}
