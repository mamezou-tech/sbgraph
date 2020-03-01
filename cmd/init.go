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
		fmt.Println("init called")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	//initCmd.PersistentFlags().String("p", "help-jp", "Name of Scrapbox project")
	initCmd.Flags().StringP("project", "p", "help-jp", "Name of Scrapbox project")
}
