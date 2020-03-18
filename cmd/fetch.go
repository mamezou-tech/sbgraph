package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "fetch all pages of the project",
	Long: `fetch all pages of the project`,
	Run: func(cmd *cobra.Command, args []string) {
		doFetch(cmd)
	},
}

func init() {
	fetchCmd.PersistentFlags().StringP("project", "p", "help-jp", "Name of Scrapbox project (required)")
	rootCmd.AddCommand(fetchCmd)
}

func doFetch(cmd *cobra.Command) {
	project, _ := cmd.PersistentFlags().GetString("project")
	fmt.Println("fetch all pages, project : ", project)
}