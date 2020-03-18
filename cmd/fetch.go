package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kondoumh/scrapbox-viz/pkg"
	"github.com/kondoumh/scrapbox-viz/pkg/fetch"
	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "fetch all pages of the project",
	Long:  `fetch all pages of the project`,
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
	count, err := fetchCount(project)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("fetch all pages, %s : %d\n", project, count)
}

func fetchCount(projectName string) (int, error) {
	url := fmt.Sprintf("%s/%s?limit=1", fetch.BaseURL, projectName)
	data, err := fetch.FetchData(url)
	if err != nil {
		return -1, err
	}
	var project pkg.Project
	err = json.Unmarshal(data, &project)
	if err != nil {
		return -1, err
	}
	return project.Count, nil
}
