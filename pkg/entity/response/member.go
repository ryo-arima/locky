package response

import "time"

// MemberResponse represents the response body for member-related operations.
// swagger:model MemberResponse
type MemberResponse struct {
	// The response code.
	//
	// required: true
	// example: "SUCCESS"
	Code string `json:"code"`
	// The response message.
	//
	// required: true
	// example: "Members retrieved successfully"
	Message string `json:"message"`
	// The list of members.
	//
	// required: true
	Members []Member `json:"members"`
}

// Member represents a member in the system.
// swagger:model Member
type Member struct {
	// The ID of the member.
	//
	// required: true
	// example: 1
	ID uint `json:"id"`
	// The UUID of the member.
	//
	// required: true
	// example: "f3b3b3b3-3b3b-3b3b-3b3b-3b3b3b3b3b3b"
	UUID string `json:"uuid"`
	// The UUID of the group.
	//
	// required: true
	// example: "f3b3b3b3-3b3b-3b3b-3b3b-3b3b3b3b3b3b"
	GroupUUID string `json:"group_uuid"`
	// The UUID of the user.
	//
	// required: true
	// example: "f3b3b3b3-3b3b-3b3b-3b3b-3b3b3b3b3b3b"
	UserUUID string `json:"user_uuid"`
	// The role of the member.
	//
	// required: true
	// example: "admin"
	Role string `json:"role"`
	// The timestamp of when the member was created.
	//
	// required: false
	// example: "2023-01-01T00:00:00Z"
	CreatedAt *time.Time `json:"created_at"`
	// The timestamp of when the member was last updated.
	//
	// required: false
	// example: "2023-01-01T00:00:00Z"
	UpdatedAt *time.Time `json:"updated_at"`
	// The timestamp of when the member was deleted.
	//
	// required: false
	// example: "2023-01-01T00:00:00Z"
	DeletedAt *time.Time `json:"deleted_at"`
}
