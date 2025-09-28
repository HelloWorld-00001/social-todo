package entity

import "errors"

var (
	// Create
	ErrCannotCreateFile    = errors.New("cannot create File")
	ErrInvalidFileCreation = errors.New("invalid File creation parameter")

	// Update
	ErrCannotUpdateFile  = errors.New("cannot update File")
	ErrInvalidFileUpdate = errors.New("invalid File update parameter")

	// Upload
	ErrCannotUploadFile  = errors.New("cannot upload File")
	ErrInvalidFileUpload = errors.New("invalid File upload parameter")

	// Delete
	ErrCannotDeleteFile = errors.New("cannot delete File")

	// General
	ErrMissingUploadFile = errors.New("missing  upload file")
	ErrCannotOpenFile    = errors.New("cannot open the uploaded file")
	ErrCannotReadFile    = errors.New("cannot read the uploaded file")
	ErrFileConflict      = errors.New("File conflict")
)
