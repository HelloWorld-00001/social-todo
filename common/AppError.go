package common

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

var (
	ErrUnauthorized    = errors.New("unauthorized access")
	ErrForbidden       = errors.New("forbidden")
	ErrConflict        = errors.New("conflict")
	ErrRequestTimeout  = errors.New("request timeout")
	ErrTooManyRequests = errors.New("too many requests")
	ErrInvalidInput    = errors.New("invalid input")
	ErrDatabase        = errors.New("database error")
)

var (
	ErrUnauthorizedMessage    = "Unauthorized access"
	ErrForbiddenMessage       = "Access denied"
	ErrConflictMessage        = "Conflict occurred"
	ErrRequestTimeoutMessage  = "Request timed out"
	ErrTooManyRequestsMessage = "Too many requests"
	ErrInvalidInputMessage    = "Invalid input provided"
	ErrDatabaseMessage        = "Internal database error"
	ErrRecordNotFoundMessage  = "Cannot find record"
)

// --- Error message formats ---
const (
	ErrCannotFindFormat   = "Cannot find %s"
	ErrCannotCreateFormat = "Cannot create %s"
	ErrCannotDeleteFormat = "Cannot delete %s"
	ErrCannotUpdateFormat = "Cannot update %s"
	ErrCannotGetFormat    = "Cannot get %s"
)

const (
	KeyUnauthorized    = "Unauthorized_Error"
	KeyForbidden       = "Forbidden_Error"
	KeyConflict        = "Conflict_Error"
	KeyRequestTimeout  = "RequestTimeout_Error"
	KeyTooManyRequests = "TooManyRequests_Error"
	KeyInvalidInput    = "InvalidInput_Error"
	KeyDatabase        = "DB_Error"
	KeyNotFound        = "NotFound_Error"
	KeyInternal        = "InternalServer_Error"
)

// Recovery is a reusable panic recovery helper.
// Call it inside defer to catch panics and log them.
func Recovery() {
	if r := recover(); r != nil {
		log.Printf("Recovered from panic: %v", r)
	}
}

type AppError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	RootErr error  `json:"-"`
	Log     string `json:"log"`
	Key     string `json:"key"`
}

func (err *AppError) Error() string {
	return err.RootError().Error()
}

func (err *AppError) RootError() error {
	var e *AppError
	if errors.As(err.RootErr, &e) {
		return e.RootError()
	}
	return err.RootErr
}

func NewFullErrorResponse(err error, message string, log string, key string, status int) *AppError {
	return &AppError{
		Message: message,
		RootErr: err,
		Log:     log,
		Key:     key,
		Status:  status,
	}
}

func NewCustomErrorResponse(err error, message string, key string) *AppError {
	return &AppError{
		Message: message,
		RootErr: err,
		Key:     key,
		Status:  http.StatusInternalServerError,
	}
}
func NewUnauthorizedError(log string) *AppError {
	return NewFullErrorResponse(ErrUnauthorized, ErrUnauthorizedMessage, log, KeyUnauthorized, http.StatusUnauthorized)
}

func NewUnauthorizedErrorCustom(err error, message string) *AppError {
	return NewFullErrorResponse(err, message, err.Error(), KeyUnauthorized, http.StatusUnauthorized)
}

func NewForbiddenError(log string) *AppError {
	return NewFullErrorResponse(ErrForbidden, ErrForbiddenMessage, log, KeyForbidden, http.StatusForbidden)
}

func NewConflictError(log string) *AppError {
	return NewFullErrorResponse(ErrConflict, ErrConflictMessage, log, KeyConflict, http.StatusConflict)
}

func NewRequestTimeoutError(log string) *AppError {
	return NewFullErrorResponse(ErrRequestTimeout, ErrRequestTimeoutMessage, log, KeyRequestTimeout, http.StatusRequestTimeout)
}

func NewTooManyRequestsError(log string) *AppError {
	return NewFullErrorResponse(ErrTooManyRequests, ErrTooManyRequestsMessage, log, KeyTooManyRequests, http.StatusTooManyRequests)
}

func NewInvalidInputError(err error) *AppError {
	return NewFullErrorResponse(err, ErrInvalidInputMessage, err.Error(), KeyInvalidInput, http.StatusBadRequest)
}

func NewDatabaseError(err error) *AppError {
	return NewFullErrorResponse(err, ErrDatabaseMessage, err.Error(), KeyDatabase, http.StatusInternalServerError)
}

func NewNotFoundErrorResponse(err error) *AppError {
	return NewFullErrorResponse(err, ErrRecordNotFoundMessage, err.Error(), KeyNotFound, http.StatusNotFound)
}

func NewBadRequestResponseWithError(err error, message string, log string) *AppError {
	return NewFullErrorResponse(err, message, log, KeyInvalidInput, http.StatusBadRequest)
}

func NewBadRequestResponse(message string) *AppError {
	return NewFullErrorResponse(errors.New(message), message, "", KeyInvalidInput, http.StatusBadRequest)
}

func NewInternalSeverErrorResponse(err error, message string, log string) *AppError {
	return NewFullErrorResponse(err, message, log, KeyInternal, http.StatusInternalServerError)
}

func NewCannotGetEntity(entity string, err error) *AppError {
	msg := fmt.Sprintf(ErrCannotGetFormat, entity)
	var appErr *AppError
	if errors.As(err, &appErr) {
		return NewFullErrorResponse(err.(*AppError).RootErr, msg, err.Error(), "Notfound"+entity, err.(*AppError).Status)
	}
	return NewFullErrorResponse(err, msg, err.Error(), "Notfound"+entity, http.StatusInternalServerError)
}

func NewCannotCreateEntity(entity string, err error) *AppError {
	msg := fmt.Sprintf(ErrCannotCreateFormat, entity)
	return NewFullErrorResponse(err, msg, err.Error(), "ErrCannotCreateFormat"+entity, http.StatusInternalServerError)
}

func NewCannotDeleteEntity(entity string, err error) *AppError {
	msg := fmt.Sprintf(ErrCannotDeleteFormat, entity)
	return NewFullErrorResponse(err, msg, err.Error(), "ErrCannotDeleteFormat"+entity, http.StatusInternalServerError)
}

func NewCannotUpdateEntity(entity string, err error) *AppError {
	msg := fmt.Sprintf(ErrCannotUpdateFormat, entity)
	return NewFullErrorResponse(err, msg, err.Error(), "ErrCannotUpdateFormat"+entity, http.StatusInternalServerError)
}

func NewInvalidUsernameOrPassword(msg string) *AppError {
	return NewFullErrorResponse(errors.New(msg), msg, "", "Error_InvalidUsernameOrPassword", http.StatusBadRequest)
}
