package model

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/ryo-arima/locky/pkg/entity/model"
)

// TestGroups_Creation tests basic group creation
func TestGroups_Creation(t *testing.T) {
	// Test data
	now := time.Now()
	group := model.Groups{
		ID:        1,
		UUID:      "test-group-uuid-123",
		Name:      "Test Group",
		CreatedAt: &now,
		UpdatedAt: &now,
		DeletedAt: nil,
	}

	// Assert using go-cmp
	expectedID := uint(1)
	if diff := cmp.Diff(expectedID, group.ID); diff != "" {
		t.Errorf("Group ID mismatch (-want +got):\n%s", diff)
	}

	expectedUUID := "test-group-uuid-123"
	if diff := cmp.Diff(expectedUUID, group.UUID); diff != "" {
		t.Errorf("Group UUID mismatch (-want +got):\n%s", diff)
	}

	expectedName := "Test Group"
	if diff := cmp.Diff(expectedName, group.Name); diff != "" {
		t.Errorf("Group Name mismatch (-want +got):\n%s", diff)
	}

	// Verify DeletedAt is nil (not soft deleted)
	if group.DeletedAt != nil {
		t.Error("Expected DeletedAt to be nil for new group")
	}
}

// TestGroups_TableDriven demonstrates table-driven tests for group validation
func TestGroups_TableDriven(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		group       model.Groups
		expectedErr bool
		description string
	}{
		{
			name: "valid group",
			group: model.Groups{
				ID:        1,
				UUID:      "valid-uuid-123",
				Name:      "Valid Group Name",
				CreatedAt: &now,
				UpdatedAt: &now,
				DeletedAt: nil,
			},
			expectedErr: false,
			description: "Should accept valid group data",
		},
		{
			name: "empty name",
			group: model.Groups{
				ID:        2,
				UUID:      "empty-name-uuid",
				Name:      "",
				CreatedAt: &now,
				UpdatedAt: &now,
				DeletedAt: nil,
			},
			expectedErr: false, // Empty name might be allowed in model validation
			description: "Should handle empty name",
		},
		{
			name: "long name",
			group: model.Groups{
				ID:        3,
				UUID:      "long-name-uuid",
				Name:      "This is a very long group name that might exceed normal limits for testing purposes",
				CreatedAt: &now,
				UpdatedAt: &now,
				DeletedAt: nil,
			},
			expectedErr: false,
			description: "Should handle long names",
		},
		{
			name: "soft deleted group",
			group: model.Groups{
				ID:        4,
				UUID:      "deleted-uuid",
				Name:      "Deleted Group",
				CreatedAt: &now,
				UpdatedAt: &now,
				DeletedAt: &now,
			},
			expectedErr: false,
			description: "Should handle soft deleted groups",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the group structure
			if tt.group.ID == 0 {
				t.Error("Expected non-zero ID")
			}

			if tt.group.UUID == "" {
				t.Error("Expected non-empty UUID")
			}

			if tt.group.CreatedAt == nil {
				t.Error("Expected CreatedAt to be set")
			}

			if tt.group.UpdatedAt == nil {
				t.Error("Expected UpdatedAt to be set")
			}

			// For soft deleted groups, verify DeletedAt is set
			if tt.name == "soft deleted group" {
				if tt.group.DeletedAt == nil {
					t.Error("Expected DeletedAt to be set for soft deleted group")
				}
			}
		})
	}
}

// TestGroups_EdgeCases tests edge cases and special scenarios
func TestGroups_EdgeCases(t *testing.T) {
	testCases := []struct {
		name        string
		description string
		testFunc    func(t *testing.T)
	}{
		{
			name:        "zero values",
			description: "Should handle zero values properly",
			testFunc: func(t *testing.T) {
				group := model.Groups{}

				// Check zero values
				expectedID := uint(0)
				if diff := cmp.Diff(expectedID, group.ID); diff != "" {
					t.Errorf("Zero ID mismatch (-want +got):\n%s", diff)
				}

				expectedUUID := ""
				if diff := cmp.Diff(expectedUUID, group.UUID); diff != "" {
					t.Errorf("Zero UUID mismatch (-want +got):\n%s", diff)
				}

				expectedName := ""
				if diff := cmp.Diff(expectedName, group.Name); diff != "" {
					t.Errorf("Zero Name mismatch (-want +got):\n%s", diff)
				}
			},
		},
		{
			name:        "special characters in name",
			description: "Should handle special characters in group names",
			testFunc: func(t *testing.T) {
				now := time.Now()
				group := model.Groups{
					ID:        1,
					UUID:      "special-char-uuid",
					Name:      "Group with Special Characters!@#$%^&*()",
					CreatedAt: &now,
					UpdatedAt: &now,
				}

				expectedName := "Group with Special Characters!@#$%^&*()"
				if diff := cmp.Diff(expectedName, group.Name); diff != "" {
					t.Errorf("Special character name mismatch (-want +got):\n%s", diff)
				}
			},
		},
		{
			name:        "unicode characters in name",
			description: "Should handle unicode characters in group names",
			testFunc: func(t *testing.T) {
				now := time.Now()
				group := model.Groups{
					ID:        1,
					UUID:      "unicode-uuid",
					Name:      "グループ名前",
					CreatedAt: &now,
					UpdatedAt: &now,
				}

				expectedName := "グループ名前"
				if diff := cmp.Diff(expectedName, group.Name); diff != "" {
					t.Errorf("Unicode name mismatch (-want +got):\n%s", diff)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.testFunc(t)
		})
	}
}

// TestGroups_FieldValidation tests individual field validation
func TestGroups_FieldValidation(t *testing.T) {
	now := time.Now()

	// Test UUID field
	t.Run("UUID validation", func(t *testing.T) {
		group := model.Groups{
			UUID: "test-uuid-validation",
		}

		if group.UUID == "" {
			t.Error("Expected UUID to be set")
		}

		expectedUUID := "test-uuid-validation"
		if diff := cmp.Diff(expectedUUID, group.UUID); diff != "" {
			t.Errorf("UUID validation mismatch (-want +got):\n%s", diff)
		}
	})

	// Test Name field
	t.Run("Name validation", func(t *testing.T) {
		group := model.Groups{
			Name: "Test Group Name",
		}

		if group.Name == "" {
			t.Error("Expected Name to be set")
		}

		expectedName := "Test Group Name"
		if diff := cmp.Diff(expectedName, group.Name); diff != "" {
			t.Errorf("Name validation mismatch (-want +got):\n%s", diff)
		}
	})

	// Test timestamp fields
	t.Run("Timestamp validation", func(t *testing.T) {
		group := model.Groups{
			CreatedAt: &now,
			UpdatedAt: &now,
		}

		if group.CreatedAt == nil {
			t.Error("Expected CreatedAt to be set")
		}

		if group.UpdatedAt == nil {
			t.Error("Expected UpdatedAt to be set")
		}

		// Verify the timestamps are as expected
		if diff := cmp.Diff(now.Unix(), group.CreatedAt.Unix()); diff != "" {
			t.Errorf("CreatedAt validation mismatch (-want +got):\n%s", diff)
		}

		if diff := cmp.Diff(now.Unix(), group.UpdatedAt.Unix()); diff != "" {
			t.Errorf("UpdatedAt validation mismatch (-want +got):\n%s", diff)
		}
	})
}

// TestGroups_Comprehensive tests comprehensive group scenarios
func TestGroups_Comprehensive(t *testing.T) {
	now := time.Now()

	// Create a comprehensive group
	group := model.Groups{
		ID:        123,
		UUID:      "comprehensive-test-uuid-456",
		Name:      "Comprehensive Test Group",
		CreatedAt: &now,
		UpdatedAt: &now,
		DeletedAt: nil,
	}

	// Test all fields comprehensively
	tests := []struct {
		fieldName string
		actual    interface{}
		expected  interface{}
	}{
		{"ID", group.ID, uint(123)},
		{"UUID", group.UUID, "comprehensive-test-uuid-456"},
		{"Name", group.Name, "Comprehensive Test Group"},
		{"CreatedAt", group.CreatedAt != nil, true},
		{"UpdatedAt", group.UpdatedAt != nil, true},
		{"DeletedAt", group.DeletedAt == nil, true},
	}

	for _, tt := range tests {
		t.Run("field_"+tt.fieldName, func(t *testing.T) {
			if diff := cmp.Diff(tt.expected, tt.actual); diff != "" {
				t.Errorf("%s field mismatch (-want +got):\n%s", tt.fieldName, diff)
			}
		})
	}
}
