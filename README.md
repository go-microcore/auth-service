# <img src="https://microcore.dev/favicon.svg" width="20" height="20" alt="Logo" /> Microcore Auth Service

[![Build](https://github.com/go-microcore/auth-service/actions/workflows/build.yml/badge.svg)](https://github.com/go-microcore/auth-service/actions/workflows/build.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-microcore/auth-service)](https://goreportcard.com/report/github.com/go-microcore/auth-service)
[![Docker Pulls](https://img.shields.io/docker/pulls/microcorelab/auth-service.svg)](https://hub.docker.com/r/microcorelab/auth-service)

Ultimate security shield for the platform. It issues and validates JWT tokens instantly, supports two-factor authentication, manages user roles and access rights, and oversees device sessions across the board. Built-in telemetry lets you see everything happening in your system, while strict security standards protect critical operations. Integrate this service, and gain full control and peace of mind over your application security.

## Features

- JWT token management
  - Issuance/renewal
  - Secure validation
  - Authorization
  - TTL configuration
  - Key generation tools
  - Management of permanent tokens
  - 2FA support
- Role-Based Access Control (RBAC)
  - HTTP rules support
- Device session management
  - List of device sessions
  - Terminate current/selected/all sessions
- Telemetry support (OpenTelemetry)

## Documentation
[https://microcore.dev/services/auth](https://microcore.dev/services/auth)


## License

This project is licensed under the terms of the [MIT License](LICENSE).

Copyright © 2026 Microcore
