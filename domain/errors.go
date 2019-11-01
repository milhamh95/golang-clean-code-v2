package domain

import (
	"github.com/friendsofgo/errors"
)

var (
	// ErrNotFound is an error message when a resource is not found
	ErrNotFound = errors.New("resource is not found")
)
