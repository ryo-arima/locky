package anonymous

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	anonymousClientBinary = "./bin/locky-anonymous"
)

// RunCommand executes anonymous client command via shell
func RunCommand(args ...string) (string, error) {
	cmd := exec.Command(anonymousClientBinary, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Set test config file path
	// Note: StartTestServer() changes cwd to project root, so we use relative path from there
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}
	testConfigPath := filepath.Join(cwd, "test", ".etc", "app.yaml")
	cmd.Env = append(os.Environ(), "CONFIG_FILE="+testConfigPath)

	err = cmd.Run()
	// Always return stdout even if command failed, as CLI outputs JSON to stdout regardless of success/error
	output := strings.TrimSpace(stdout.String())

	if err != nil {
		// If no stdout but have stderr, return error with stderr
		if output == "" && stderr.Len() > 0 {
			return "", fmt.Errorf("command failed: %w, stderr: %s", err, stderr.String())
		}
		// Otherwise return stdout (which contains JSON error response) with no error
	}

	return output, nil
}

// Create commands

func CreateUser(name, email, password string) (string, error) {
	return RunCommand("create", "user", "--name", name, "--email", email, "--password", password, "-o", "json")
}

// Common (Auth) commands

func Login(email, password string) (string, error) {
	return RunCommand("common", "login", "--email", email, "--password", password, "-o", "json")
}

func ValidateToken(token string) (string, error) {
	return RunCommand("common", "validate", "--access-token", token, "-o", "json")
}

func GetUserInfo(token string) (string, error) {
	return RunCommand("common", "userinfo", "--access-token", token, "-o", "json")
}

func RefreshToken(token string) (string, error) {
	return RunCommand("common", "refresh", "--refresh-token", token, "-o", "json")
}

func Logout(token string) (string, error) {
	return RunCommand("common", "logout", "--access-token", token, "-o", "json")
}
