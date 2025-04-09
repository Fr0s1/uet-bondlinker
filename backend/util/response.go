
package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is a standardized API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Count   int64       `json:"count,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// RespondWithSuccess sends a successful response with data
func RespondWithSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// RespondWithError sends an error response
func RespondWithError(c *gin.Context, statusCode int, errorMessage string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error:   errorMessage,
	})
}

// RespondWithCount sends a successful response with data and count
func RespondWithCount(c *gin.Context, statusCode int, message string, data interface{}, count int64) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
		Count:   count,
	})
}

// RespondWithPagination sends a successful response with pagination metadata
func RespondWithPagination(c *gin.Context, statusCode int, message string, data interface{}, meta interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// RespondWithNoContent sends a 204 No Content response
func RespondWithNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// ErrorMessages provides standardized error messages
var ErrorMessages = struct {
	NotAuthenticated     string
	NotAuthorized        string
	InvalidRequest       string
	ResourceNotFound     string
	InternalServerError  string
	Conflict             string
	ValidationFailed     string
	DatabaseError        string
	AlreadyExists        string
	InvalidToken         string
	ExpiredToken         string
	InvalidCredentials   string
	InvalidUserID        string
	InvalidResourceID    string
}{
	NotAuthenticated:     "Not authenticated",
	NotAuthorized:        "Not authorized to perform this action",
	InvalidRequest:       "Invalid request parameters",
	ResourceNotFound:     "Resource not found",
	InternalServerError:  "Internal server error",
	Conflict:             "Resource conflict",
	ValidationFailed:     "Validation failed",
	DatabaseError:        "Database error",
	AlreadyExists:        "Resource already exists",
	InvalidToken:         "Invalid or expired token",
	ExpiredToken:         "Token has expired",
	InvalidCredentials:   "Invalid credentials",
	InvalidUserID:        "Invalid user ID",
	InvalidResourceID:    "Invalid resource ID",
}
