package application

import "errors"

var (
	ErrBadRequest          = errors.New("bad request")
	ErrNotFound            = errors.New("page not found")
	ErrMethodNotAllowed    = errors.New("method not allowed")
	ErrInternalServerError = errors.New("internal server error")
	ErrUnknown             = errors.New("unknown error")
)
