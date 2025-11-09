package app

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	appClientBinary = "./bin/locky-app"
)

// RunCommand executes app client command via shell
func RunCommand(args ...string) (string, error) {
	cmd := exec.Command(appClientBinary, args...)
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

func CreateGroup(name string) (string, error) {
	return RunCommand("create", "group", "--name", name, "-o", "json")
}

func CreateMember(userID, groupID string) (string, error) {
	return RunCommand("create", "member", "--user-id", userID, "--group-id", groupID, "-o", "json")
}

// Get commands

func GetUser(id string) (string, error) {
	return RunCommand("get", "user", "--id", id, "-o", "json")
}

func GetUserList() (string, error) {
	return RunCommand("get", "users", "-o", "json")
}

func GetGroup(id string) (string, error) {
	return RunCommand("get", "group", "--id", id, "-o", "json")
}

func GetGroupList() (string, error) {
	return RunCommand("get", "groups", "-o", "json")
}

func GetMember(id string) (string, error) {
	return RunCommand("get", "member", "--id", id, "-o", "json")
}

func GetMemberList() (string, error) {
	return RunCommand("get", "member", "--list", "-o", "json")
}

func GetRole(id string) (string, error) {
	return RunCommand("get", "role", "--id", id, "-o", "json")
}

func GetRoleList() (string, error) {
	return RunCommand("get", "roles", "-o", "json")
}

// Update commands

func UpdateUser(id, name, email string) (string, error) {
	return RunCommand("update", "user", "--id", id, "--name", name, "--email", email, "-o", "json")
}

func UpdateGroup(id, name string) (string, error) {
	return RunCommand("update", "group", "--id", id, "--name", name, "-o", "json")
}

func UpdateMember(id, userID, groupID string) (string, error) {
	return RunCommand("update", "member", "--id", id, "--user-id", userID, "--group-id", groupID, "-o", "json")
}

// Delete commands

func DeleteUser(id string) (string, error) {
	return RunCommand("delete", "user", "--id", id, "-o", "json")
}

func DeleteGroup(id string) (string, error) {
	return RunCommand("delete", "group", "--id", id, "-o", "json")
}

func DeleteMember(id string) (string, error) {
	return RunCommand("delete", "member", "--id", id, "-o", "json")
}
