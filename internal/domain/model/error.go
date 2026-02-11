package model

import "errors"

var (
	ErrForbidden      = errors.New("access forbidden")
	ErrNotFound       = errors.New("not found")
	ErrBadRequest     = errors.New("bad request")
	ErrInvalidExtFile = errors.New("invalid file extension")
)
