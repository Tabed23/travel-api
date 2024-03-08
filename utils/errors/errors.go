package errors

import "errors"
var (
	ErrUnAuthorized = errors.New("unAuthorized user")
	ErrInValidRole = errors.New("invalid role")
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserNotFound = errors.New("user not found")
	ErrBadRequest = errors.New("bad request")
	ErrInternalServerError = errors.New("internal server error")
)