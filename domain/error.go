package domain

import (
	"fmt"
	"net/http"

	"github.com/friendsofgo/errors"
)

var (
	// ErrNotFound is an error message when a resource is not found
	ErrNotFound = errors.New("resource is not found")

	// ErrNotModified is thrown to the client when the cached copy of a partifulcar file is up to date with the server
	ErrNotModified = errors.New("")
)

// ConstraintError representes a custom error for a constraint things
type ConstraintError string

func (e ConstraintError) Error() string {
	return string(e)
}

// ConstraintErrorf constructs ConstraintError with formatted message
func ConstraintErrorf(format string, a ...interface{}) ConstraintError {
	return ConstraintError(fmt.Sprintf(format, a...))
}

// ErrorFromResponseStatusCode generates error based on the status code from *http.Response.
// For example, it will generate ErrNotFound when given status code is 404
func ErrorFromResponseStatusCode(code int, message string) (err error) {
	switch code {
	case http.StatusNotFound:
		err = ErrNotFound
	case http.StatusBadRequest:
		err = ConstraintErrorf(message)
	case http.StatusNotModified:
		err = ErrNotModified
	default:
		err = fmt.Errorf(message)
	}

	return
}
