package middleware

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ErrorMiddleware returns an error with response http status code
func ErrorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return echo.NewHTTPError(http.StatusOK, "ok")
			}

			fmt.Println("========  ========")
			fmt.Printf("%+v\n", err)
			fmt.Println("=================")

			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
	}
}
