package errors

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

// ErrorData - Core error structure
type ErrorData struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Error   error  `json:"-"`
}

// ErrorResponse - Client-facing error response
type ErrorResponse struct {
	ErrorData
	Variables  V                `json:"variables,omitempty"`
	Validation ValidationErrors `json:"validation,omitempty"`
	StatusCode int              `json:"-"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("code=%d, message=%s", e.Code, e.Message)
}

// Error codes
const (
	CodeBadRequest       = 4000
	CodeValidationError  = 4001
	CodeInvalidPackSizes = 4002
	CodeMalformedRequest = 4003
	CodeInternalServer   = 5000
	CodeTooManyRequests  = 4290
)

// Predefined errors
var (
	ErrBadRequest = ErrorData{
		Message: "Bad request",
		Code:    CodeBadRequest,
	}

	ErrValidation = ErrorData{
		Message: "Validation error",
		Code:    CodeValidationError,
	}

	ErrInvalidPackSizes = ErrorData{
		Message: "Invalid pack sizes provided",
		Code:    CodeInvalidPackSizes,
	}

	ErrMalformedRequest = ErrorData{
		Message: "Malformed request",
		Code:    CodeMalformedRequest,
	}

	ErrInternalServer = ErrorData{
		Message: "Internal server error",
		Code:    CodeInternalServer,
	}

	ErrTooManyRequests = ErrorData{
		Message: "Too many requests",
		Code:    CodeTooManyRequests,
	}
)

// NewBadRequestError creates a 400 error
func NewBadRequestError(base ErrorData, cause error, vars V) *ErrorResponse {
	resp := &ErrorResponse{
		ErrorData: ErrorData{
			Message: base.Message,
			Code:    base.Code,
			Error:   cause,
		},
		Variables:  vars,
		StatusCode: http.StatusBadRequest,
	}

	// Add validation details if available
	if ve, ok := cause.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			resp.Validation = append(resp.Validation, FieldError{
				Field:   fe.Field(),
				Message: fe.Tag(),
			})
		}
	}

	return resp
}

// NewInternalServerError creates a 500 error
func NewInternalServerError(base ErrorData, cause error, vars V) *ErrorResponse {
	return &ErrorResponse{
		ErrorData: ErrorData{
			Message: base.Message,
			Code:    base.Code,
			Error:   cause,
		},
		Variables:  vars,
		StatusCode: http.StatusInternalServerError,
	}
}

// NewRateLimitError creates a 429 error
func NewRateLimitError(retryAfter string) *ErrorResponse {
	return &ErrorResponse{
		ErrorData: ErrorData{
			Message: ErrTooManyRequests.Message,
			Code:    ErrTooManyRequests.Code,
		},
		Variables: V{
			"retry_after": retryAfter,
		},
		StatusCode: http.StatusTooManyRequests,
	}
}

// Abort sends the error response and stops the request chain
func (e *ErrorResponse) Abort(c *gin.Context) {
	c.AbortWithStatusJSON(e.StatusCode, gin.H{
		"code":       e.Code,
		"message":    e.Message,
		"details":    e.Variables,
		"validation": e.Validation,
	})
}
