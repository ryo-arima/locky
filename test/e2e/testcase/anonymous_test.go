package testcase

import (
	"encoding/json"
	"testing"

	"github.com/ryo-arima/locky/test/e2e/client/anonymous"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnonymousUserRegistration(t *testing.T) {
	t.Run("Create User", func(t *testing.T) {
		output, err := anonymous.CreateUser("user1", "user1@locky.local", "User1Password123!")
		require.NoError(t, err, "User registration should succeed")

		var result map[string]interface{}
		err = json.Unmarshal([]byte(output), &result)
		if err != nil {
			t.Logf("Failed to parse response: %s", output)
		}
		require.NoError(t, err, "Should parse JSON output")

		// Response structure: {code, message, users: [{id, name, email, uuid}]}
		if users, ok := result["users"].([]interface{}); ok && len(users) > 0 {
			if user, ok := users[0].(map[string]interface{}); ok {
				assert.NotEmpty(t, user["id"], "User ID should be returned")
				assert.Equal(t, "user1", user["name"])
				assert.Equal(t, "user1@locky.local", user["email"])
			}
		}
	})
}

func TestAuthenticationFlow(t *testing.T) {
	var accessToken string
	var refreshToken string
	testEmail := "user2@locky.local"
	testPassword := "User2Password123!"

	// First create a user
	t.Run("Create User for Auth", func(t *testing.T) {
		output, err := anonymous.CreateUser("user2", testEmail, testPassword)
		require.NoError(t, err, "User creation should succeed")
		assert.NotEmpty(t, output, "Should return user data")
	})

	// Login
	t.Run("Login", func(t *testing.T) {
		output, err := anonymous.Login(testEmail, testPassword)
		require.NoError(t, err, "Login should succeed")

		var result map[string]interface{}
		err = json.Unmarshal([]byte(output), &result)
		if err != nil {
			t.Logf("Failed to parse login response: %s", output)
		}
		require.NoError(t, err, "Should parse JSON output")

		// Extract tokens from response: {"token_pair": {"access_token": "...", "refresh_token": "..."}}
		if tokenPair, ok := result["token_pair"].(map[string]interface{}); ok {
			if at, ok := tokenPair["access_token"].(string); ok {
				accessToken = at
			}
			if rt, ok := tokenPair["refresh_token"].(string); ok {
				refreshToken = rt
			}
		}
		// Fallback for old format
		if accessToken == "" {
			if tokenData, ok := result["token"].(map[string]interface{}); ok {
				accessToken = tokenData["access_token"].(string)
				if rt, ok := tokenData["refresh_token"].(string); ok {
					refreshToken = rt
				}
			} else if at, ok := result["access_token"].(string); ok {
				accessToken = at
			}
		}
		assert.NotEmpty(t, accessToken, "Access token should be returned")
		assert.NotEmpty(t, refreshToken, "Refresh token should be returned")
	})

	// Validate token
	t.Run("Validate Token", func(t *testing.T) {
		output, err := anonymous.ValidateToken(accessToken)
		require.NoError(t, err, "Token validation should succeed")

		var result map[string]interface{}
		err = json.Unmarshal([]byte(output), &result)
		if err != nil {
			t.Logf("Failed to parse ValidateToken response: %s", output)
		}
		require.NoError(t, err, "Should parse JSON output")

		// Response structure: {"code": "SUCCESS", "message": "Token is valid", "data": {...}}
		assert.Equal(t, "SUCCESS", result["code"], "Should return success code")
		assert.NotNil(t, result["data"], "Should return token data")
	})

	// Get user info from token
	t.Run("Get User Info", func(t *testing.T) {
		output, err := anonymous.GetUserInfo(accessToken)
		require.NoError(t, err, "Get user info should succeed")

		var result map[string]interface{}
		err = json.Unmarshal([]byte(output), &result)
		if err != nil {
			t.Logf("Failed to parse GetUserInfo response: %s", output)
		}
		require.NoError(t, err, "Should parse JSON output")

		// Response structure: {"code": "SUCCESS", "message": "...", "commons": [{id, uuid, ...}]}
		assert.Equal(t, "SUCCESS", result["code"], "Should return success code")
		if commons, ok := result["commons"].([]interface{}); ok && len(commons) > 0 {
			if common, ok := commons[0].(map[string]interface{}); ok {
				assert.NotEmpty(t, common["id"], "User ID should be present")
				assert.NotEmpty(t, common["uuid"], "User UUID should be present")
			}
		}
	})

	// Refresh token
	var newAccessToken string
	t.Run("Refresh Token", func(t *testing.T) {
		output, err := anonymous.RefreshToken(refreshToken)
		require.NoError(t, err, "Token refresh should succeed")

		var result map[string]interface{}
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err, "Should parse JSON output")

		// Extract new tokens from response: {"token_pair": {"access_token": "...", "refresh_token": "..."}}
		if tokenPair, ok := result["token_pair"].(map[string]interface{}); ok {
			if at, ok := tokenPair["access_token"].(string); ok {
				newAccessToken = at
			}
		}
		// Fallback for old format
		if newAccessToken == "" {
			if tokenData, ok := result["token"].(map[string]interface{}); ok {
				newAccessToken = tokenData["access_token"].(string)
			} else if at, ok := result["access_token"].(string); ok {
				newAccessToken = at
			}
		}

		assert.NotEmpty(t, newAccessToken, "New access token should be returned")
		assert.NotEqual(t, accessToken, newAccessToken, "New access token should be different")
	})

	// Logout
	t.Run("Logout", func(t *testing.T) {
		output, err := anonymous.Logout(newAccessToken)
		require.NoError(t, err, "Logout should succeed")
		assert.NotEmpty(t, output, "Should return logout confirmation")
	})

	// Validate token after logout (should fail)
	t.Run("Validate Token After Logout", func(t *testing.T) {
		output, err := anonymous.ValidateToken(newAccessToken)
		// CLI returns JSON even for errors, so check the response code instead
		require.NoError(t, err, "Command should execute")

		var result map[string]interface{}
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err, "Should parse JSON output")

		// Should return error code indicating token is revoked
		code, ok := result["code"].(string)
		require.True(t, ok, "Response should have code field")
		assert.NotEqual(t, "SUCCESS", code, "Token validation should fail after logout")
	})
}
