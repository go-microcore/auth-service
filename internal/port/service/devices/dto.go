// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package devices

type (
	// Results

	// DeviceResult represents a device with its active session.
	DeviceResult struct {
		ID      string
		Session SessionResult
	}

	// SessionResult represents the details of a device session.
	SessionResult struct {
		IssuedAt       string
		Location       string
		IP             string
		UserAgent      string
		OsFullName     string
		OsName         string
		OsVersion      string
		Platform       string
		Model          string
		BrowserName    string
		BrowserVersion string
		EngineName     string
		EngineVersion  string
	}
)
