package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Set current project",
	Long: LongUsage(
		`Set current project

		  sbraph project -p <project name>

		project will be saved to config.
		`),
	Run: func(cmd *cobra.Command, args []string) {
		doProject(cmd)
	},
}

func init() {
	projectCmd.PersistentFlags().StringP("project", "p", "help-jp", "Name of Scrapbox project")
	rootCmd.AddCommand(projectCmd)
}

func doProject(cmd *cobra.Command) {
	projectName, _ := cmd.PersistentFlags().GetString("project")
	fmt.Printf("set current project : %s\n", projectName)
	viper.Set("currentproject", projectName)
	SaveConfig()
}
