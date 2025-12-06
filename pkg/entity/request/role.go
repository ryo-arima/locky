package request

// RolePermissionItem: permission element (resource, action)
// swagger:model RolePermissionItem
type RolePermissionItem struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

// RolePermission: role creation/update request body
// swagger:model RolePermission
type RolePermission struct {
	Role        string               `json:"role"`        // role name
	Permissions []RolePermissionItem `json:"permissions"` // permissions list
}
