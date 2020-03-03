package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// graphCmd represents the graph command
var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "generate graph structure",
	Long: `generate graph structure of pages and authors (as dot file)`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("graph called")
	},
}

func init() {
	rootCmd.AddCommand(graphCmd)
}
