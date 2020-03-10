package cmd

import (
	"fmt"
	"os"

	"github.com/kondoumh/scrapbox-viz/pkg/util"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize directory",
	Long:  `initalize directory for Scrapbox project`,
	Run: func(cmd *cobra.Command, args []string) {
		project, _ := cmd.PersistentFlags().GetString("project")
		fmt.Println("init called, project : ", project)
		if err := filesystem.CreateDir(config.WorkDir); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	initCmd.PersistentFlags().StringP("project", "p", "help-jp", "Name of Scrapbox project (required)")
	rootCmd.AddCommand(initCmd)
}
