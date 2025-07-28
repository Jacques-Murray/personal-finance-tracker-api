package responses

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// ValidationErrorResponse represents a standardized validation error response
type ValidationErrorResponse struct {
	Error  string                 `json:"error"`
	Fields []ValidationFieldError `json:"fields"`
}

// ValidationFieldError represents details for a specific field validation error
type ValidationFieldError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Message string `json:"message,omitempty"`
}
