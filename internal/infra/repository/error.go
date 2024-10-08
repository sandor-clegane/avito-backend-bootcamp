package repository

import "errors"

var (
	ErrNotFound            = errors.New("entity not found")
	ErrConstraintViolation = errors.New("db constraint violation")
	ErrAlreadyExists       = errors.New("unique constraint violation")
)
