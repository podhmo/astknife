package action

import "github.com/pkg/errors"

var (
	// ErrNoEffect :
	ErrNoEffect = errors.New("no effect")
	// ErrNotFound :
	ErrNotFound = errors.New("not found")
)
