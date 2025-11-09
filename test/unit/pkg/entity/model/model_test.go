package model_test

import (
	"testing"
	"time"

	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/stretchr/testify/assert"
)

func TestUsers_Structure(t *testing.T) {
	now := time.Now()
	user := model.Users{
		ID:        1,
		UUID:      "test-uuid-123",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		Name:      "Test User",
		CreatedAt: &now,
		UpdatedAt: &now,
		DeletedAt: nil,
	}

	assert.Equal(t, uint(1), user.ID)
	assert.Equal(t, "test-uuid-123", user.UUID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "hashedpassword", user.Password)
	assert.Equal(t, "Test User", user.Name)
	assert.NotNil(t, user.CreatedAt)
	assert.NotNil(t, user.UpdatedAt)
	assert.Nil(t, user.DeletedAt)
}

func TestGroups_Structure(t *testing.T) {
	now := time.Now()
	group := model.Groups{
		ID:        1,
		UUID:      "group-uuid-123",
		Name:      "Test Group",
		CreatedAt: &now,
		UpdatedAt: &now,
		DeletedAt: nil,
	}

	assert.Equal(t, uint(1), group.ID)
	assert.Equal(t, "group-uuid-123", group.UUID)
	assert.Equal(t, "Test Group", group.Name)
	assert.NotNil(t, group.CreatedAt)
	assert.NotNil(t, group.UpdatedAt)
	assert.Nil(t, group.DeletedAt)
}

func TestMembers_Structure(t *testing.T) {
	now := time.Now()
	member := model.Members{
		ID:        1,
		UUID:      "member-uuid-123",
		UserID:    10,
		GroupID:   20,
		CreatedAt: &now,
		UpdatedAt: &now,
		DeletedAt: nil,
	}

	assert.Equal(t, uint(1), member.ID)
	assert.Equal(t, "member-uuid-123", member.UUID)
	assert.Equal(t, uint(10), member.UserID)
	assert.Equal(t, uint(20), member.GroupID)
	assert.NotNil(t, member.CreatedAt)
	assert.NotNil(t, member.UpdatedAt)
	assert.Nil(t, member.DeletedAt)
}

func TestRoles_Structure(t *testing.T) {
	now := time.Now()
	role := model.Roles{
		ID:        1,
		UUID:      "role-uuid-123",
		Name:      "Admin",
		CreatedAt: &now,
		UpdatedAt: &now,
		DeletedAt: nil,
	}

	assert.Equal(t, uint(1), role.ID)
	assert.Equal(t, "role-uuid-123", role.UUID)
	assert.Equal(t, "Admin", role.Name)
	assert.NotNil(t, role.CreatedAt)
	assert.NotNil(t, role.UpdatedAt)
	assert.Nil(t, role.DeletedAt)
}
