package database

import "github.com/pkg/errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrDuplicateKey = errors.New("already exists")
	ErrCtxDone      = errors.New("ctx done error")
)
