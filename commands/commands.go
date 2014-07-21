package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "goplanet",
	Short: `planet style RSS aggregator`,
	Long:  `Provides a planet style RSS aggregator`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("goplanet runs")
	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
