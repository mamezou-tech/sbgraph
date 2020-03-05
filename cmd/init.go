package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize directory",
	Long: `initalize directory for Scrapbox project`,
	Run: func(cmd *cobra.Command, args []string) {
		project, _ := cmd.PersistentFlags().GetString("project")
		fmt.Println("init called, project : ", project)
	},
}

func init() {
	initCmd.PersistentFlags().StringP("project", "p", "help-jp", "Name of Scrapbox project (required)")
	rootCmd.AddCommand(initCmd)
}
