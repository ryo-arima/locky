package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ryo-arima/locky/pkg/config"
)

func main() {
	ctx := context.Background()

	fmt.Println("=== Locky Configuration Test ===")
	fmt.Println()

	// Test 1: File-based configuration (default)
	fmt.Println("Test 1: File-based configuration")
	os.Setenv("USE_SECRETSMANAGER", "false")
	cfg1 := config.NewBaseConfigWithContext(ctx)
	fmt.Printf("  MySQL Host: %s\n", cfg1.YamlConfig.MySQL.Host)
	fmt.Printf("  MySQL User: %s\n", cfg1.YamlConfig.MySQL.User)
	fmt.Printf("  MySQL DB: %s\n", cfg1.YamlConfig.MySQL.Db)
	fmt.Printf("  Redis Host: %s\n", cfg1.YamlConfig.Redis.Host)
	fmt.Printf("  Server JWT Secret: %s\n\n", cfg1.YamlConfig.Application.Server.JWTSecret)

	// Test 2: Secrets Manager with LocalStack
	fmt.Println("Test 2: Secrets Manager (LocalStack)")
	os.Setenv("USE_SECRETSMANAGER", "true")
	os.Setenv("USE_LOCALSTACK", "true")
	os.Setenv("SECRET_ID", "locky/config/app")
	os.Setenv("AWS_ENDPOINT_URL", "http://localhost:4566")
	os.Setenv("AWS_REGION", "us-east-1")

	cfg2 := config.NewBaseConfigWithContext(ctx)
	fmt.Printf("  MySQL Host: %s\n", cfg2.YamlConfig.MySQL.Host)
	fmt.Printf("  MySQL User: %s\n", cfg2.YamlConfig.MySQL.User)
	fmt.Printf("  MySQL DB: %s\n", cfg2.YamlConfig.MySQL.Db)
	fmt.Printf("  Redis Host: %s\n", cfg2.YamlConfig.Redis.Host)
	fmt.Printf("  Server JWT Secret: %s\n\n", cfg2.YamlConfig.Application.Server.JWTSecret)

	fmt.Println("=== All tests completed successfully ===")
}
