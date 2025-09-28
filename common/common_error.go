package common

import "errors"

var (
	ErrInvalidInputID         = errors.New("invalid id parameter")
	ErrInvalidInputPagination = errors.New("invalid pagination parameter")
	ErrInvalidInputFilter     = errors.New("invalid filter parameter")
	ErrInvalidEntity          = errors.New("invalid entity parameter")
	ErrUnhandleError          = errors.New("internal server error")
	ErrInvalidInput           = errors.New("invalid input")
	ErrDatabaseError          = errors.New("database error")
)
