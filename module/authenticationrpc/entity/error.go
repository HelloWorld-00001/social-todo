package entity

import "errors"

var (
	// Create
	ErrCannotCreateAccount    = errors.New("cannot create account")
	ErrInvalidAccountCreation = errors.New("invalid account creation parameter")

	// Update
	ErrCannotUpdateAccount  = errors.New("cannot update account")
	ErrInvalidAccountUpdate = errors.New("invalid account update parameter")

	// Disable/Activate
	ErrCannotDisableAccount  = errors.New("cannot disable account")
	ErrCannotActivateAccount = errors.New("cannot activate account")

	// Login
	ErrInvalidLoginCredential = errors.New("invalid username or password")
	ErrCannotGenerateToken    = errors.New("cannot generate authentication token")

	// General
	ErrCannotFindAccount = errors.New("cannot find account")
	ErrAccountConflict   = errors.New("account already exists")
)
