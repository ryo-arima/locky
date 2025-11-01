package request

// RolePermissionItem: permission element (resource, action)
// swagger:model RolePermissionItem
type RolePermissionItem struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

// RolePermissionRequest: role creation/update request body
// swagger:model RolePermissionRequest
type RolePermissionRequest struct {
	Role        string               `json:"role"`        // role name
	Permissions []RolePermissionItem `json:"permissions"` // permissions list
}
