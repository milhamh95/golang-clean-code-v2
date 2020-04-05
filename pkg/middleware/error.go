package middleware

import (
	"context"
	"net/http"

	"github.com/friendsofgo/errors"
	"github.com/labstack/echo/v4"

	"github.com/milhamhidayat/golang-clean-code-v2/domain"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// ErrorMiddleware returns an error with response http status code
func ErrorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			err = errors.Cause(err)

			if _, ok := err.(domain.ConstraintError); ok {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			switch err {
			case context.DeadlineExceeded, context.Canceled:
				return echo.NewHTTPError(http.StatusRequestTimeout, err.Error())
			case domain.ErrNotFound:
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			case domain.ErrNotModified:
				return c.NoContent(http.StatusNotModified)
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
}
