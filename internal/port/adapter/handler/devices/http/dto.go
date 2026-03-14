// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

package http

type (
	// Responses

	// DeviceResponse represents the API response for a device, including its session details.
	DeviceResponse struct {
		ID      string          `json:"id"`
		Session SessionResponse `json:"session"`
	}

	// SessionResponse contains details about a user session, such as location, device, OS,
	// and browser information.
	SessionResponse struct {
		IssuedAt       string `json:"issuedAt"`
		Location       string `json:"location"`
		IP             string `json:"ip"`
		UserAgent      string `json:"userAgent"`
		OsFullName     string `json:"osFullName"`
		OsName         string `json:"osName"`
		OsVersion      string `json:"osVersion"`
		Platform       string `json:"platform"`
		Model          string `json:"model"`
		BrowserName    string `json:"browserName"`
		BrowserVersion string `json:"browserVersion"`
		EngineName     string `json:"engineName"`
		EngineVersion  string `json:"engineVersion"`
	}
)
