package testcase

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/ryo-arima/locky/test/e2e/client/anonymous"
	"github.com/ryo-arima/locky/test/e2e/client/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppGroupCRUD(t *testing.T) {
	// Login first to get access token
	t.Run("Setup App User Login", func(t *testing.T) {
		// Create app user if not exists
		_, _ = anonymous.CreateUser("appuser", "appuser@locky.local", "AppUser123!")
		
		output, err := anonymous.Login("appuser@locky.local", "AppUser123!")
		require.NoError(t, err, "Login should succeed")

		var result map[string]interface{}
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err, "Should parse login response")

		// Extract token from response: {"token_pair": {"access_token": "..."}}
		var accessToken string
		if tokenPair, ok := result["token_pair"].(map[string]interface{}); ok {
			if at, ok := tokenPair["access_token"].(string); ok {
				accessToken = at
			}
		}
		require.NotEmpty(t, accessToken, "Access token should be returned")
		os.Setenv("LOCKY_ACCESS_TOKEN", accessToken)
	})

	var groupID string

	// Create group
	t.Run("Create Group", func(t *testing.T) {
		output, err := app.CreateGroup("app-test-group")
		t.Logf("CreateGroup output: %q, err: %v", output, err)
		require.NoError(t, err, "Create group should succeed")

		var result map[string]interface{}
		err = json.Unmarshal([]byte(output), &result)
		if err != nil {
			t.Logf("Failed to parse CreateGroup response: %s", output)
		}
		require.NoError(t, err, "Should parse JSON output")

		// GroupResponse structure: {code, message, groups: [{id, uuid, name}]}
		if groups, ok := result["groups"].([]interface{}); ok && len(groups) > 0 {
			if group, ok := groups[0].(map[string]interface{}); ok {
				if id, ok := group["id"].(float64); ok {
					groupID = fmt.Sprintf("%.0f", id)
				}
			}
		}
		assert.NotEmpty(t, groupID, "Group ID should be returned")
	})

	// List groups to verify creation
	t.Run("List Groups", func(t *testing.T) {
		output, err := app.GetGroupList()
		require.NoError(t, err, "List groups should succeed")

		var result map[string]interface{}
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err, "Should parse JSON output")

		if groups, ok := result["groups"].([]interface{}); ok {
			assert.NotEmpty(t, groups, "Group list should not be empty")
		}
	})

	// Update group
	t.Run("Update Group", func(t *testing.T) {
		output, err := app.UpdateGroup(groupID, "updated-app-group")
		require.NoError(t, err, "Update group should succeed")

		var result map[string]interface{}
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err, "Should parse JSON output")

		if groups, ok := result["groups"].([]interface{}); ok && len(groups) > 0 {
			if group, ok := groups[0].(map[string]interface{}); ok {
				assert.Equal(t, "updated-app-group", group["name"])
			}
		}
	})

	// Delete group
	t.Run("Delete Group", func(t *testing.T) {
		output, err := app.DeleteGroup(groupID)
		require.NoError(t, err, "Delete group should succeed")
		assert.NotEmpty(t, output, "Delete should return confirmation")
	})
}

func TestAppUserOperations(t *testing.T) {
	// App users can only read their own user info and other users (read-only)
	t.Run("List Users", func(t *testing.T) {
		output, err := app.GetUserList()
		require.NoError(t, err, "List users should succeed")

		var result []map[string]interface{}
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err, "Should parse JSON array")

		// May be empty or contain users depending on test state
		t.Logf("Found %d users", len(result))
	})
}

func TestAppRoleReadOnly(t *testing.T) {
	// App users can only read roles
	t.Run("List Roles", func(t *testing.T) {
		output, err := app.GetRoleList()
		require.NoError(t, err, "List roles should succeed")

		var result []map[string]interface{}
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err, "Should parse JSON array")

		t.Logf("Found %d roles", len(result))
	})
}
