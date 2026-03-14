// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package logout

type (
	// Data

	// DeviceData represents the data required to log out a specific device for a user.
	DeviceData struct {
		User   uint
		Device string
	}

	// AllData represents the data required to log out all devices for a user.
	AllData struct {
		User uint
	}
)
