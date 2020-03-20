package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kondoumh/scrapbox-viz/pkg/file"

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
	projectName, _ := cmd.PersistentFlags().GetString("project")
	project, err := fetchIndex(projectName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("fetch all pages, %s : %d\n", project.Name, project.Count)
	err = fetchPageList(project)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func fetchIndex(projectName string) (pkg.Project, error) {
	url := fmt.Sprintf("%s/%s?limit=1", fetch.BaseURL, projectName)
	data, err := fetch.FetchData(url)
	var project pkg.Project
	if err != nil {
		return project, err
	}
	err = json.Unmarshal(data, &project)
	if err != nil {
		return project, err
	}
	return project, nil
}

func fetchPageList(project pkg.Project) error {
	pages := []pkg.Page{}
	for skip := 0; skip < project.Count; skip += fetch.Limit {
		url := fmt.Sprintf("%s/%s?skip=%d&limit=%d&sort=updated", fetch.BaseURL, project.Name, skip, fetch.Limit)
		data, err := fetch.FetchData(url)
		if err != nil {
			return err
		}
		var proj pkg.Project
		err = json.Unmarshal(data, &proj)
		for _, page := range proj.Pages {
			pages = append(pages, page)
		}
	}
	project.Pages = pages
	data, _ := json.Marshal(project)
	if err := file.WriteBytes(data, project.Name+".json", config.WorkDir); err != nil {
		return err
	}
	return nil
}
