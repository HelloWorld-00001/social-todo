package entity

import "errors"

var (
	// Create
	ErrCannotCreateUser    = errors.New("cannot create user")
	ErrInvalidUserCreation = errors.New("invalid user creation parameter")

	// Update
	ErrCannotUpdateUser  = errors.New("cannot update user")
	ErrInvalidUserUpdate = errors.New("invalid user update parameter")

	// Disable/Activate
	ErrCannotDisableUser  = errors.New("cannot disable user")
	ErrCannotActivateUser = errors.New("cannot activate user")

	// Login
	ErrInvalidLoginCredential = errors.New("invalid username or password")
	ErrCannotGenerateToken    = errors.New("cannot generate authentication token")

	// General
	ErrCannotFindUser = errors.New("cannot find user")
	ErrUserConflict   = errors.New("user already exists")
)
