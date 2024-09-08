package errors

import "errors"

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidDatabaseType = errors.New("invalid database type")
)
