package utils

// ErrorResponse represents the structure of error responses returned by the API.
// This is used for Swagger documentation.
type ErrorResponse struct {
	Status  string      `json:"status" example:"error"`
	Message string      `json:"message" example:"An error message explaining what went wrong"`
	Details interface{} `json:"details,omitempty"`
}

// SuccessResponse represents the structure of success responses returned by the API.
// This is used for Swagger documentation.
type SuccessResponse struct {
	Status  string      `json:"status" example:"success"`
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
}
