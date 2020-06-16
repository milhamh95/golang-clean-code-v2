package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	departmentHandler "github.com/milhamhidayat/golang-clean-code-v2/department/delivery/http"
	"github.com/milhamhidayat/golang-clean-code-v2/pkg/middleware"
)

const address = ":8500"

var serverCmd = &cobra.Command{
	Use:   "http",
	Short: "Start http server",
	Run: func(cmd *cobra.Command, args []string) {
		e := echo.New()
		e.Use(middleware.ErrorMiddleware())

		e.GET("ping", func(c echo.Context) error {
			return c.JSON(http.StatusOK, "pong")
		})

		departmentHandler.AddDepartmentHandler(e, departmentService)

		errCh := make(chan error)

		go func(ch chan error) {
			log.Info().Msgf("Starting HTTP server at: %s", address)
			errCh <- e.Start(address)
		}(errCh)

		go func(ch chan error) {
			errCh <- http.ListenAndServe(":6060", nil)
		}(errCh)

		for {
			log.Fatal().Err(<-errCh)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
