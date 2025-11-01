package response

// RoleResponse: role operation response
// swagger:model RoleResponse
type RoleResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Roles   interface{} `json:"roles"`
	Detail  interface{} `json:"detail,omitempty"`
}
