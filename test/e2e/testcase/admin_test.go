package testcase

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/ryo-arima/locky/test/e2e/client/admin"
	"github.com/ryo-arima/locky/test/e2e/client/anonymous"
	"github.com/ryo-arima/locky/test/e2e/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Start test server
	if err := server.StartTestServer(); err != nil {
		panic("Failed to start test server: " + err.Error())
	}
	defer server.StopTestServer()

	// Initialize database schema
	if err := server.InitializeDatabase(); err != nil {
		panic("Failed to initialize database: " + err.Error())
	}

	// Run tests
	m.Run()
}

func TestAdminUserCRUD(t *testing.T) {
	// Bootstrap initial schema first (recreates tables)
	t.Run("Bootstrap User Schema", func(t *testing.T) {
		_, _ = admin.BootstrapUser("", "", "")
		// Bootstrap just creates/recreates the table, ignore errors
	})

	// Create admin user via anonymous API for authentication
	t.Run("Setup Admin User", func(t *testing.T) {
		output, err := anonymous.CreateUser("admin", "admin@locky.local", "AdminPassword123!")
		t.Logf("Create admin user output: %s, error: %v", output, err)
		// User might already exist, that's ok - just log it
		if err == nil {
			var result map[string]interface{}
			if jsonErr := json.Unmarshal([]byte(output), &result); jsonErr == nil {
				t.Logf("Admin user created successfully: %+v", result)
			}
		}
	})

	// Login to get access token
	var token string
	t.Run("Login as Admin", func(t *testing.T) {
		output, err := anonymous.Login("admin@locky.local", "AdminPassword123!")
		require.NoError(t, err, "Login should succeed")

		var result map[string]interface{}
		err = json.Unmarshal([]byte(output), &result)
		if err != nil {
			t.Logf("Failed to parse login response: %s", output)
		}
		require.NoError(t, err, "Should parse login response")

		// Extract token from response: {"token_pair": {"access_token": "..."}}
		if tokenPair, ok := result["token_pair"].(map[string]interface{}); ok {
			if accessToken, ok := tokenPair["access_token"].(string); ok {
				token = accessToken
				os.Setenv("LOCKY_ACCESS_TOKEN", token)
			}
		}
		// Fallback for old format
		if token == "" {
			if tokenData, ok := result["token"].(map[string]interface{}); ok {
				token = tokenData["access_token"].(string)
				os.Setenv("LOCKY_ACCESS_TOKEN", token)
			}
		}
		require.NotEmpty(t, token, "Access token should be returned")
	})

	// Create user
	var userID uint
	t.Run("Create User", func(t *testing.T) {
		output, err := admin.CreateUser("test1", "test1@locky.local", "TestPassword123!")
		if err != nil {
			t.Logf("Create user error output: %s", output)
		}
		require.NoError(t, err, "Create user should succeed")

		var result map[string]interface{}
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err, "Should parse JSON output")

		// Response structure: {code, message, users: [{id, name, email, ...}]}
		if users, ok := result["users"].([]interface{}); ok && len(users) > 0 {
			if user, ok := users[0].(map[string]interface{}); ok {
				if id, ok := user["id"].(float64); ok {
					userID = uint(id)
				}
			}
		}
		require.NotZero(t, userID, "User ID should be returned")
	})

	// Get user
	t.Run("Get User", func(t *testing.T) {
		output, err := admin.GetUser(fmt.Sprintf("%d", userID))
		require.NoError(t, err, "Get user should succeed")

		var result map[string]interface{}
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err, "Should parse JSON output")

		assert.Equal(t, float64(userID), result["id"])
		assert.Equal(t, "test1", result["name"])
	})

	// Update user
	t.Run("Update User", func(t *testing.T) {
		output, err := admin.UpdateUser(fmt.Sprintf("%d", userID), "test1-updated", "test1-updated@locky.local")
		require.NoError(t, err, "Update user should succeed")

		var result map[string]interface{}
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err, "Should parse JSON output")

		assert.Equal(t, "test1-updated", result["name"])
		assert.Equal(t, "test1-updated@locky.local", result["email"])
	})

	// List users
	t.Run("List Users", func(t *testing.T) {
		output, err := admin.GetUserList()
		require.NoError(t, err, "List users should succeed")

		var result []map[string]interface{}
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err, "Should parse JSON array")

		assert.NotEmpty(t, result, "User list should not be empty")
	})

	// Delete user
	t.Run("Delete User", func(t *testing.T) {
		output, err := admin.DeleteUser(fmt.Sprintf("%d", userID))
		require.NoError(t, err, "Delete user should succeed")
		assert.NotEmpty(t, output, "Delete should return confirmation")
	})
}

func TestAdminGroupCRUD(t *testing.T) {
	var groupID string

	// Create group
	t.Run("Create Group", func(t *testing.T) {
		output, err := admin.CreateGroup("test-group", "Test group description")
		require.NoError(t, err, "Create group should succeed")

		var result map[string]interface{}
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err, "Should parse JSON output")

		groupID = result["id"].(string)
		assert.NotEmpty(t, groupID, "Group ID should be returned")
		assert.Equal(t, "test-group", result["name"])
	})

	// Get group
	t.Run("Get Group", func(t *testing.T) {
		output, err := admin.GetGroup(groupID)
		require.NoError(t, err, "Get group should succeed")

		var result map[string]interface{}
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err, "Should parse JSON output")

		assert.Equal(t, groupID, result["id"])
		assert.Equal(t, "test-group", result["name"])
	})

	// Update group
	t.Run("Update Group", func(t *testing.T) {
		output, err := admin.UpdateGroup(groupID, "updated-group", "Updated description")
		require.NoError(t, err, "Update group should succeed")

		var result map[string]interface{}
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err, "Should parse JSON output")

		assert.Equal(t, "updated-group", result["name"])
	})

	// List groups
	t.Run("List Groups", func(t *testing.T) {
		output, err := admin.GetGroupList()
		require.NoError(t, err, "List groups should succeed")

		var result []map[string]interface{}
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err, "Should parse JSON array")

		assert.NotEmpty(t, result, "Group list should not be empty")
	})

	// Delete group
	t.Run("Delete Group", func(t *testing.T) {
		output, err := admin.DeleteGroup(groupID)
		require.NoError(t, err, "Delete group should succeed")
		assert.NotEmpty(t, output, "Delete should return confirmation")
	})
}

func TestAdminRoleCRUD(t *testing.T) {
	var roleID string

	// Create role
	t.Run("Create Role", func(t *testing.T) {
		output, err := admin.CreateRole("test-role", "Test role description")
		require.NoError(t, err, "Create role should succeed")

		var result map[string]interface{}
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err, "Should parse JSON output")

		roleID = result["id"].(string)
		assert.NotEmpty(t, roleID, "Role ID should be returned")
		assert.Equal(t, "test-role", result["name"])
	})

	// Get role
	t.Run("Get Role", func(t *testing.T) {
		output, err := admin.GetRole(roleID)
		require.NoError(t, err, "Get role should succeed")

		var result map[string]interface{}
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err, "Should parse JSON output")

		assert.Equal(t, roleID, result["id"])
		assert.Equal(t, "test-role", result["name"])
	})

	// Update role
	t.Run("Update Role", func(t *testing.T) {
		output, err := admin.UpdateRole(roleID, "updated-role", "Updated description")
		require.NoError(t, err, "Update role should succeed")

		var result map[string]interface{}
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err, "Should parse JSON output")

		assert.Equal(t, "updated-role", result["name"])
	})

	// List roles
	t.Run("List Roles", func(t *testing.T) {
		output, err := admin.GetRoleList()
		require.NoError(t, err, "List roles should succeed")

		var result []map[string]interface{}
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err, "Should parse JSON array")

		assert.NotEmpty(t, result, "Role list should not be empty")
	})

	// Delete role
	t.Run("Delete Role", func(t *testing.T) {
		output, err := admin.DeleteRole(roleID)
		require.NoError(t, err, "Delete role should succeed")
		assert.NotEmpty(t, output, "Delete should return confirmation")
	})
}
