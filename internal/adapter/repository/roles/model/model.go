// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

// Package model defines database models for roles repository adapter.
package model

import (
	"time"

	"go.microcore.dev/auth-service/internal/adapter/repository/rules/model"
)

type (
	// AuthRole represents a role stored in the database.
	AuthRole struct {
		ID          string               `gorm:"primarykey"`
		Name        string               `gorm:"uniqueIndex;not null"`
		Description string               `gorm:"not null"`
		SystemFlag  bool                 `gorm:"not null"`
		ServiceFlag bool                 `gorm:"not null"`
		Created     time.Time            `gorm:"not null"`
		Updated     time.Time            `gorm:"not null"`
		HTTPRules   []model.AuthHTTPRule `gorm:"foreignKey:RoleID;references:ID"`
	}
)
