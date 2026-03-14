// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

// Package model defines database models for auth repository adapter.
package model

import (
	"time"

	"github.com/lib/pq"
)

type (
	// AuthStaticToken represents a static access token stored in the database.
	AuthStaticToken struct {
		ID          string         `gorm:"primarykey"`
		Token       string         `gorm:"uniqueIndex;not null;size:64"`
		UserID      uint           `gorm:"not null"`
		Device      string         `gorm:"not null"`
		Roles       pq.StringArray `gorm:"type:text[];not null"`
		Description string         `gorm:"not null"`
		Created     time.Time      `gorm:"not null"`
	}
)
