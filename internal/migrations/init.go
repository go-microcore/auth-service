// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// AuthInit creates initial tables and indexes.
func AuthInit() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "auth_init",
		Migrate: func(tx *gorm.DB) error {
			// Create tables
			if err := tx.Exec(`
				CREATE TABLE IF NOT EXISTS auth_roles (
					id TEXT PRIMARY KEY,
					name TEXT NOT NULL UNIQUE,
					description TEXT NOT NULL,
					system_flag BOOLEAN NOT NULL,
					service_flag BOOLEAN NOT NULL,
					created TIMESTAMPTZ NOT NULL,
					updated TIMESTAMPTZ NOT NULL
				);
			`).Error; err != nil {
				return err
			}

			if err := tx.Exec(`
				CREATE TABLE IF NOT EXISTS auth_http_rules (
					id SERIAL PRIMARY KEY,
					role_id TEXT NOT NULL,
					path TEXT NOT NULL,
					methods TEXT[] NOT NULL,
					mfa BOOLEAN DEFAULT TRUE,
					created TIMESTAMPTZ NOT NULL,
					updated TIMESTAMPTZ NOT NULL,
					CONSTRAINT fk_role FOREIGN KEY (role_id)
						REFERENCES auth_roles(id)
						ON UPDATE CASCADE
						ON DELETE CASCADE
				);
			`).Error; err != nil {
				return err
			}

			if err := tx.Exec(`
				CREATE TABLE IF NOT EXISTS auth_static_tokens (
					id TEXT PRIMARY KEY,
					token CHAR(64) NOT NULL UNIQUE,
					user_id BIGINT NOT NULL,
					device TEXT NOT NULL,
					roles TEXT[] NOT NULL,
					description TEXT,
					created TIMESTAMPTZ NOT NULL
				);
			`).Error; err != nil {
				return err
			}

			// Create indexes

			if err := tx.Exec(`
				CREATE INDEX IF NOT EXISTS idx_auth_http_rules_role_id 
				ON auth_http_rules(role_id);
			`).Error; err != nil {
				return err
			}

			if err := tx.Exec(`
				CREATE UNIQUE INDEX idx_auth_static_tokens_token 
				ON auth_static_tokens(token);
			`).Error; err != nil {
				return err
			}

			return nil
		},

		Rollback: func(tx *gorm.DB) error {
			// Drop indexes
			if err := tx.Exec(`
				DROP INDEX IF EXISTS idx_auth_http_rules_role_id;
			`).Error; err != nil {
				return err
			}

			if err := tx.Exec(`
				DROP INDEX IF EXISTS idx_auth_static_tokens_token;
			`).Error; err != nil {
				return err
			}

			// Drop tables

			if err := tx.Exec(`
				DROP TABLE IF EXISTS auth_static_tokens;
			`).Error; err != nil {
				return err
			}

			if err := tx.Exec(`
				DROP TABLE IF EXISTS auth_http_rules;
			`).Error; err != nil {
				return err
			}

			if err := tx.Exec(`
				DROP TABLE IF EXISTS auth_roles;
			`).Error; err != nil {
				return err
			}

			return nil
		},
	}
}
