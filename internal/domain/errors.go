package domain

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrConflict         = errors.New("already exists")
	ErrUnknownAlgorithm = errors.New("unknown algorithm")
	ErrInvalidID        = errors.New("invalid device id")
)
