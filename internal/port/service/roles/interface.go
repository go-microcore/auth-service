// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package roles

import (
	"context"
)

type (
	// Service defines the interface for the roles service.
	Service interface {
		CreateRole(
			ctx context.Context,
			data CreateRoleData,
		) (*CreateRoleResult, error)

		FilterRoles(
			ctx context.Context,
			data FilterRolesData,
		) ([]FilterRolesResult, error)

		DeleteRole(
			ctx context.Context,
			id string,
		) error

		UpdateRole(
			ctx context.Context,
			id string,
			data map[string]any,
		) error

		AuthorizeHTTPRoles(
			ctx context.Context,
			data AuthorizeHTTPRolesData,
		) (*AuthorizeHTTPRolesResult, error)
	}
)
