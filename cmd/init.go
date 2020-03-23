package cmd

import (
	"fmt"
	"os"

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
	if err := file.CreateDir(config.WorkDir); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := fetchProject(project); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func fetchProject(project string) error {
	url := fmt.Sprintf("%s/%s?limit=1", api.BaseURL, project)
	data, err := api.Fetch(url)
	if err != nil {
		return err
	}

	if err := file.WriteBytes(data, project+".json", config.WorkDir); err != nil {
		return err
	}
	return nil
}
