package response

// CountResponse generic count response
// swagger:model CountResponse
// Simple response that returns only count in addition to code and message
// Example: {"code":"SUCCESS","message":"Count retrieved","count":123}
type CountResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Count   int64  `json:"count"`
}
