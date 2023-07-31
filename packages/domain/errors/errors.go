package errors

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserDoesNotExist  = errors.New("user does not exist")
	ErrWebsiteNotFound   = errors.New("website not found")
)
