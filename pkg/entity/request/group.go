package request

import "time"

// GroupRequest represents the request body for group-related operations.
// swagger:model GroupRequest
type GroupRequest struct {
	// The ID of the group.
	//
	// required: false
	// example: 1
	ID uint `json:"id"`
	// The UUID of the group.
	//
	// required: false
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
