package response_test

import (
	"testing"

	"github.com/ryo-arima/locky/pkg/entity/response"
	"github.com/stretchr/testify/assert"
)

func TestUserResponse_Structure(t *testing.T) {
	resp := response.UserResponse{
		Code:    "SUCCESS",
		Message: "Operation successful",
		Data: []map[string]interface{}{
			{
				"uuid":  "user-uuid-123",
				"email": "test@example.com",
				"name":  "Test User",
			},
		},
	}

	assert.Equal(t, "SUCCESS", resp.Code)
	assert.Equal(t, "Operation successful", resp.Message)
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, "user-uuid-123", resp.Data[0]["uuid"])
}

func TestGroupResponse_Structure(t *testing.T) {
	resp := response.GroupResponse{
		Code:    "SUCCESS",
		Message: "Group operation successful",
		Data: []map[string]interface{}{
			{
				"uuid": "group-uuid-123",
				"name": "Test Group",
			},
		},
	}

	assert.Equal(t, "SUCCESS", resp.Code)
	assert.Equal(t, "Group operation successful", resp.Message)
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, "group-uuid-123", resp.Data[0]["uuid"])
}

func TestMemberResponse_Structure(t *testing.T) {
	resp := response.MemberResponse{
		Code:    "SUCCESS",
		Message: "Member operation successful",
		Data: []map[string]interface{}{
			{
				"uuid":       "member-uuid-123",
				"user_uuid":  "user-uuid-456",
				"group_uuid": "group-uuid-789",
			},
		},
	}

	assert.Equal(t, "SUCCESS", resp.Code)
	assert.Equal(t, "Member operation successful", resp.Message)
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, "member-uuid-123", resp.Data[0]["uuid"])
}

func TestRoleResponse_Structure(t *testing.T) {
	resp := response.RoleResponse{
		Code:    "SUCCESS",
		Message: "Role operation successful",
		Data: []map[string]interface{}{
			{
				"uuid": "role-uuid-123",
				"name": "Admin",
			},
		},
	}

	assert.Equal(t, "SUCCESS", resp.Code)
	assert.Equal(t, "Role operation successful", resp.Message)
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, "role-uuid-123", resp.Data[0]["uuid"])
}

func TestLoginResponse_Structure(t *testing.T) {
	resp := response.LoginResponse{
		Code:    "SUCCESS",
		Message: "Login successful",
		Token:   "jwt-token-abc123",
	}

	assert.Equal(t, "SUCCESS", resp.Code)
	assert.Equal(t, "Login successful", resp.Message)
	assert.Equal(t, "jwt-token-abc123", resp.Token)
}

func TestErrorResponse_Structure(t *testing.T) {
	resp := response.UserResponse{
		Code:    "ERROR",
		Message: "An error occurred",
		Data:    nil,
	}

	assert.Equal(t, "ERROR", resp.Code)
	assert.Equal(t, "An error occurred", resp.Message)
	assert.Nil(t, resp.Data)
}
