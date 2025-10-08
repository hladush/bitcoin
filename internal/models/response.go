package models

// APIResponse represents a standard API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// ErrorResponse creates a standardized error response
func ErrorResponse(message string) APIResponse {
	return APIResponse{
		Success: false,
		Error:   message,
	}
}

// SuccessResponse creates a standardized success response
func SuccessResponse(data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Data:    data,
	}
}

// MessageResponse creates a standardized message response
func MessageResponse(message string) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
	}
}
