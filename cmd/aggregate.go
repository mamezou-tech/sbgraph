package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// aggregateCmd represents the aggregate command
var aggregateCmd = &cobra.Command{
	Use:   "aggregate",
	Short: "aggregate project activities",
	Long: `aggregate project activities`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("aggregate called")
	},
}

func init() {
	rootCmd.AddCommand(aggregateCmd)
}
