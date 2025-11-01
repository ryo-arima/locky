package response

import "time"

// GroupResponse represents the response body for group-related operations.
// swagger:model GroupResponse
type GroupResponse struct {
	// The response code.
	//
	// required: true
	// example: "SUCCESS"
	Code string `json:"code"`
	// The response message.
	//
	// required: true
	// example: "Groups retrieved successfully"
	Message string `json:"message"`
	// The list of groups.
	//
	// required: true
	Groups []Group `json:"groups"`
}

// Group represents a group in the system.
// swagger:model Group
type Group struct {
	// The ID of the group.
	//
	// required: true
	// example: 1
	ID uint `json:"id"`
	// The UUID of the group.
	//
	// required: true
	// example: "f3b3b3b3-3b3b-3b3b-3b3b-3b3b3b3b3b3b"
	UUID string `json:"uuid"`
	// The name of the group.
	//
	// required: true
	// example: "My Group"
	Name string `json:"name"`
	// The timestamp of when the group was created.
	//
	// required: false
	// example: "2023-01-01T00:00:00Z"
	CreatedAt *time.Time `json:"created_at"`
	// The timestamp of when the group was last updated.
	//
	// required: false
	// example: "2023-01-01T00:00:00Z"
	UpdatedAt *time.Time `json:"updated_at"`
	// The timestamp of when the group was deleted.
	//
	// required: false
	// example: "2023-01-01T00:00:00Z"
	DeletedAt *time.Time `json:"deleted_at"`
}
