package errors

import (
	"fmt"
	"kn-assignment/internal/constant"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a structured error response
type ErrorResponse struct {
	Code    constant.ErrorCode `json:"code"`
	Message string             `json:"message"`
}

// NewErrorResponse creates a new ErrorResponse
func NewErrorResponse(code constant.ErrorCode, message string) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
	}
}

// HandleError sends a structured error response
func HandleError(c *gin.Context, code constant.ErrorCode) {
	statusCode := mapErrorCodeToHTTPStatus(code)
	c.JSON(statusCode, NewErrorResponse(code, code.String()))
	c.Abort()
}

// CustomError represents a custom error with a code and message
type CustomError struct {
	Code    constant.ErrorCode
	Message string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

// NewCustomError creates a new CustomError
func NewCustomError(code constant.ErrorCode) error {
	return &CustomError{
		Code:    code,
		Message: code.String(),
	}
}

func NewCustomErrorWithMessage(code constant.ErrorCode, message string) error {
	return &CustomError{
		Code:    code,
		Message: message,
	}
}

// mapErrorCodeToHTTPStatus maps custom error codes to HTTP status codes
func mapErrorCodeToHTTPStatus(code constant.ErrorCode) int {
	switch code {
	case constant.ErrCodeInvalidRequest:
		return http.StatusBadRequest
	case constant.ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case constant.ErrCodeForbidden:
		return http.StatusForbidden
	case constant.ErrCodeNotFound:
		return http.StatusNotFound
	case constant.ErrCodeConflict:
		return http.StatusConflict
	case constant.ErrCodeInternalServer:
		fallthrough
	default:
		return http.StatusInternalServerError
	}
}
