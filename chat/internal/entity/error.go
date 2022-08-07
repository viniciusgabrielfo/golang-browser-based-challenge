package entity

import "errors"

var (
	ErrNotFoundEntity  = errors.New("not found entity")
	ErrPassworNotMatch = errors.New("password doesn't match")
)
