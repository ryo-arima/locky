package response

import "time"

// UserResponse represents the response body for user-related operations.
// swagger:model UserResponse
type UserResponse struct {
	// The response code.
	//
	// required: true
	// example: "SUCCESS"
	Code string `json:"code"`
	// The response message.
	//
	// required: true
	// example: "Users retrieved successfully"
	Message string `json:"message"`
	// The list of users.
	//
	// required: true
	Users []User `json:"users"`
}

// User represents a user in the system.
// swagger:model User
type User struct {
	// The ID of the user.
	//
	// required: true
	// example: 1
	ID uint `json:"id"`
	// The UUID of the user.
	//
	// required: true
	// example: "f3b3b3b3-3b3b-3b3b-3b3b-3b3b3b3b3b3b"
	UUID string `json:"uuid"`
	// The email of the user.
	//
	// required: true
	// example: "jhon.doe@example.com"
	Email string `json:"email"`
	// The name of the user.
	//
	// required: true
	// example: "John Doe"
	Name string `json:"name"`
	// The timestamp of when the user was created.
	//
	// required: false
	// example: "2023-01-01T00:00:00Z"
	CreatedAt *time.Time `json:"created_at"`
	// The timestamp of when the user was last updated.
	//
	// required: false
	// example: "2023-01-01T00:00:00Z"
	UpdatedAt *time.Time `json:"updated_at"`
	// The timestamp of when the user was deleted.
	//
	// required: false
	// example: "2023-01-01T00:00:00Z"
	DeletedAt *time.Time `json:"deleted_at"`
}
