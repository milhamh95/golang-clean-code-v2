package middleware

import (
	"fmt"
	"net/http"

	"github.com/friendsofgo/errors"
	"github.com/labstack/echo/v4"

	"github.com/milhamhidayat/golang-clean-code-v2/domain"
)

// ErrorMiddleware returns an error with response http status code
func ErrorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			fmt.Println("========  ========")
			fmt.Printf("%+v\n", err)
			fmt.Println("=================")

			if errors.Is(err, domain.ErrNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, domain.ErrNotFound)
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
}
