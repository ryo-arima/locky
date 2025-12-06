package response

import "time"

// Resources represents the response body for resource-related operations.
// swagger:model Resources
type Resources struct {
	// The response code.
	//
	// required: true
	// example: "SUCCESS"
	Code string `json:"code"`
	// The response message.
	//
	// required: true
	// example: "Resources retrieved successfully"
	Message string `json:"message"`
	// The list of resources.
	//
	// required: true
	Resources []Resource `json:"resources"`
}

// Resource represents a resource in the system.
// swagger:model Resource
type Resource struct {
	// The ID of the resource.
	//
	// required: true
	// example: 1
	ID uint `json:"id"`
	// The UUID of the resource.
	//
	// required: true
	// example: "f3b3b3b3-3b3b-3b3b-3b3b-3b3b3b3b3b3b"
	UUID string `json:"uuid"`
	// The name of the resource.
	//
	// required: true
	// example: "My Resource"
	Name string `json:"name"`
	// The type of the resource.
	//
	// required: false
	// example: "compute"
	Type string `json:"type"`
	// The UUID of the owning group.
	//
	// required: true
	// example: "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	GroupUUID string `json:"group_uuid"`
	// The timestamp of when the resource was created.
	//
	// required: true
	// example: "2023-01-01T00:00:00Z"
	CreatedAt *time.Time `json:"created_at"`
	// The timestamp of when the resource was last updated.
	//
	// required: true
	// example: "2023-01-02T00:00:00Z"
	UpdatedAt *time.Time `json:"updated_at"`
	// The timestamp of when the resource was deleted.
	//
	// required: false
	// example: "2023-01-03T00:00:00Z"
	DeletedAt *time.Time `json:"deleted_at"`
}
