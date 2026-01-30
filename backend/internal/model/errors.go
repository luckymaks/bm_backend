package model

import "github.com/cockroachdb/errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrValidation    = errors.New("validation error")
)
