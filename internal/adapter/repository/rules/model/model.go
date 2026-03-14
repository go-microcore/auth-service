// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

// Package model defines database models for rules repository adapter.
package model

import (
	"time"

	"github.com/lib/pq"
)

type (
	// AuthHTTPRule represents a HTTP rule stored in the database.
	AuthHTTPRule struct {
		ID      uint           `gorm:"primarykey"`
		RoleID  string         `gorm:"not null;index"`
		Path    string         `gorm:"not null"`
		Methods pq.StringArray `gorm:"type:text[];not null"`
		Mfa     bool           `gorm:"not null"`
		Created time.Time      `gorm:"not null"`
		Updated time.Time      `gorm:"not null"`
	}
)
