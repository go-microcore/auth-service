// Copyright (c) 2026 Microcore (https://microcore.dev)
// Licensed under the MIT License. See LICENSE file for details.

// Package main is the entry point for the auth-service CLI application.
package main

import (
	"os"

	"go.microcore.dev/auth-service/cmd/cli"
	"go.microcore.dev/framework/shutdown"
)

func main() {
	os.Exit(shutdown.Exit(cli.Run()))
}
