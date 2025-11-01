package request

import "time"

// UserRequest represents the request body for user-related operations.
// swagger:model UserRequest
type UserRequest struct {
	// The ID of the user.
	//
	// required: false
	// example: 1
	ID uint `json:"id"`
	// The UUID of the user.
	//
	// required: false
	// example: "f3b3b3b3-3b3b-3b3b-3b3b-3b3b3b3b3b3b"
	UUID string `json:"uuid"`
	// The name of the user.
	//
	// required: true
	// example: "John Doe"
	Name string `json:"name"`
	// The email of the user.
	//
	// required: true
	// example: "jhon.doe@example.com"
	Email string `json:"email"`
	// The password of the user.
	//
	// required: true
	// example: "password"
	Password string `json:"password"`
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
