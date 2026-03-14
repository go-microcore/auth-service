// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
)

// Get returns all migrations.
func Get() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		AuthInit(),
	}
}
