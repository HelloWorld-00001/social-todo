package common

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ErrorHandler is a structured error for consistent handling
type ErrorHandler struct {
	RootErr  error  `json:"-"` // deepest error (DB, network, etc.)
	CauseErr error  `json:"-"` // higher-level cause
	Message  string `json:"message"`
	Log      string `json:"log"`
	Code     string `json:"code"`   // machine-friendly code
	Status   int    `json:"status"` // HTTP status
}

// Error implements error
func (e *ErrorHandler) Error() string {
	return e.Message
}

// Unwrap lets errors.Is / errors.As traverse the chain
func (e *ErrorHandler) Unwrap() error {
	if e.CauseErr != nil {
		return e.CauseErr
	}
	return e.RootErr
}

// ===== Fluent modifiers =====

// WithError attaches an intermediate cause (like domain/entity error)
func (e *ErrorHandler) WithError(err error) *ErrorHandler {
	newErr := *e
	newErr.CauseErr = err
	if err != nil {
		newErr.Log = fmt.Sprintf("%s | cause: %v", e.Message, err)
	}
	return &newErr
}

// WithRootCause attaches the deepest root error (like DB error)
func (e *ErrorHandler) WithRootCause(err error) *ErrorHandler {
	newErr := *e
	newErr.RootErr = err
	if err != nil {
		newErr.Log = fmt.Sprintf("%s | root: %v", e.Message, err)
	}
	return &newErr
}

// RootError extracts the deepest error
func RootError(err error) error {
	for {
		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			return err
		}
		err = unwrapped
	}
}

// ===== Common Error Codes =====

var (
	BadRequest = &ErrorHandler{
		Message: "bad request",
		Code:    "BadRequest",
		Status:  http.StatusBadRequest,
	}

	NotFound = &ErrorHandler{
		Message: "resource not found",
		Code:    "NotFound",
		Status:  http.StatusNotFound,
	}

	InternalServerError = &ErrorHandler{
		Message: "internal server error",
		Code:    "InternalServerError",
		Status:  http.StatusInternalServerError,
	}

	Unauthorized = &ErrorHandler{
		Message: "unauthorized",
		Code:    "Unauthorized",
		Status:  http.StatusUnauthorized,
	}

	DatabaseError = &ErrorHandler{
		Message: "database error",
		Code:    "database_error",
		Status:  http.StatusInternalServerError,
	}

	Forbidden = &ErrorHandler{
		Message: "forbidden",
		Code:    "Forbidden",
		Status:  http.StatusForbidden,
	}

	Conflict = &ErrorHandler{
		Message: "resource is conflict",
		Code:    "Conflict",
		Status:  http.StatusConflict,
	}
)

// RespondError handles returned error from business and writes JSON response.
func RespondError(c *gin.Context, err error) {
	var e *ErrorHandler
	if errors.As(err, &e) {
		c.JSON(e.Status, gin.H{
			"error": gin.H{
				"code":    e.Code,
				"message": e.Message,
			},
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"error": gin.H{
			"code":    "InternalServerError",
			"message": "something went wrong",
		},
	})
}
