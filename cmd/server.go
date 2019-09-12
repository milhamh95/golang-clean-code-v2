package cmd

import (
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const address = ":8500"

var serverCmd = &cobra.Command{
	Use:   "http",
	Short: "Start http server",
	Run: func(cmd *cobra.Command, args []string) {
		e := echo.New()

		e.GET("ping", func(c echo.Context) error {
			return c.JSON(http.StatusOK, "pong")
		})

		log.Info("Starting HTTP server at ", address)
		err := e.Start(address)
		if err != nil {
			log.Fatalf("Failed to start server: %s", err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
