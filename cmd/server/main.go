// Package main provides the Locky API server.
//
// # Locky API Server
//
// Locky is a JWT-based authentication and user management system with the following features:
//   - User authentication and authorization using JWT tokens
//   - User, Group, and Member management
//   - OpenStack Keystone-style API design
//   - REST API with multiple access levels (public, internal, private)
//
// # API Structure
//
// The API follows OpenStack Keystone design patterns:
//   - /v1/share/common/auth/* - Authentication endpoints (token management)
//   - /v1/public/* - Public API (no authentication required)
//   - /v1/internal/* - Internal API (authentication required)
//   - /v1/private/* - Private API (admin operations)
//
// # Authentication
//
// The system uses JWT tokens with the following endpoints:
//   - POST /v1/share/common/auth/tokens - Issue token (login)
//   - DELETE /v1/share/common/auth/tokens - Revoke token (logout)
//   - GET /v1/share/common/auth/tokens/validate - Validate token
//   - POST /v1/share/common/auth/tokens/refresh - Refresh token
//   - GET /v1/share/common/auth/tokens/user - Get user info from token
//
// # Configuration
//
// The server reads configuration with the following priority:
//  1. AWS Secrets Manager (if USE_SECRETSMANAGER=true and SECRET_ID is set)
//  2. Local file specified by CONFIG_FILE environment variable
//  3. Default: etc/app.yaml
//
// Environment Variables:
//   - CONFIG_FILE: Path to configuration file (default: etc/app.yaml)
//   - USE_SECRETSMANAGER: Set to "true" to use AWS Secrets Manager
//   - SECRET_ID: AWS Secrets Manager secret ID
//   - CONFIG_SOURCE: "secretsmanager" or "localfile" (alternative to USE_SECRETSMANAGER)
//
// Configuration includes:
//   - Database connection settings (MySQL/TiDB)
//   - Redis connection settings
//   - JWT secret key
//   - Admin user email list
//
// # Usage
//
//	go run cmd/server/main.go
//
// The server will start on http://localhost:8080 by default.
//
// Schemes: http, https
// Host: localhost:8080
// BasePath: /v1
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Security:
// - bearerAuth: []
//
// SecurityDefinitions:
// bearerAuth:
//
//	type: apiKey
//	name: Authorization
//	in: header
//	description: "JWT token. Format: 'Bearer {token}'"
//
// swagger:meta
package main

import (
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/logger"
	"github.com/ryo-arima/locky/pkg/server"
	"github.com/ryo-arima/locky/pkg/server/middleware"
)

func main() {
	// Set logger factory
	config.SetLoggerFactory(middleware.NewLogger)

	conf := config.NewBaseConfig()

	// Initialize global logger
	logger.Initialize(conf.YamlConfig.Logger)

	// サーバー起動時にのみ DB 接続
	_ = conf.ConnectDB()
	// Redis は router 初期化時に repository.NewRedisClient を利用 (既存実装)
	server.Main(*conf)
}
