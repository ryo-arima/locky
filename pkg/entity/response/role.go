package response

// Roles: role operation response
// swagger:model Roles
type Roles struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Roles   interface{} `json:"roles"`
	Detail  interface{} `json:"detail,omitempty"`
}
