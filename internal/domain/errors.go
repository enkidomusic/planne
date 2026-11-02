package domain

import "errors"

var (
	ErrInvalidInput = errors.New("missing required task fields")
	ErrInternal     = errors.New("internal server error")
)
