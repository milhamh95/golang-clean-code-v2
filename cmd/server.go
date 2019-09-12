package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "http",
	Short: "Start http server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("++++++++ print ++++++++")
		fmt.Println("start server")
		fmt.Println("+++++++++++++++++")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
