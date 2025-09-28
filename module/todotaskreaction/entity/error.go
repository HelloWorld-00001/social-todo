package entity

import "errors"

var (
	// Create
	ErrCannotCreateReaction    = errors.New("cannot create reaction")
	ErrInvalidReactionCreation = errors.New("invalid reaction creation parameter")

	// Update
	ErrCannotUpdateReaction  = errors.New("cannot update reaction")
	ErrInvalidReactionUpdate = errors.New("invalid reaction update parameter")

	// Delete
	ErrCannotDeleteReaction = errors.New("cannot delete reaction")

	// General
	ErrCannotFindReaction = errors.New("cannot find reaction")
	ErrReactionConflict   = errors.New("reaction conflict")
)
