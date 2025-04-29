package errors

// V is a generic map for additional error variables
type V map[string]interface{}

// FieldError represents validation error details
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors collects multiple field errors
type ValidationErrors []FieldError
