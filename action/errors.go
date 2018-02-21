package action

import "github.com/pkg/errors"

var (
	// ErrReplacementNotFound :
	ErrReplacementNotFound = errors.New("replacement not found")
	// ErrTargetNotFound :
	ErrTargetNotFound = errors.New("target not found")
)

// IsNoEffect :
func IsNoEffect(err error) bool {
	switch errors.Cause(err) {
	case ErrReplacementNotFound, ErrTargetNotFound:
		return true
	default:
		return false
	}
}
