package cmd

import (
	"fmt"

	"github.com/kondoumh/scrapbox-viz/pkg/api"
	"github.com/kondoumh/scrapbox-viz/pkg/file"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize directory",
	Long:  `initalize directory for Scrapbox project`,
	Run: func(cmd *cobra.Command, args []string) {
		doInit(cmd)
	},
}

func init() {
	initCmd.PersistentFlags().StringP("project", "p", "help-jp", "Name of Scrapbox project (required)")
	rootCmd.AddCommand(initCmd)
}

func doInit(cmd *cobra.Command) {
	project, _ := cmd.PersistentFlags().GetString("project")
	fmt.Println("init called, project : ", project)
	err := file.CreateDir(config.WorkDir)
	CheckErr(err)
	err = fetchProject(project)
	CheckErr(err)
}

func fetchProject(project string) error {
	data, err := api.FetchIndex(project)
	if err != nil {
		return err
	}

	if err := file.WriteBytes(data, project+".json", config.WorkDir); err != nil {
		return err
	}
	return nil
}
