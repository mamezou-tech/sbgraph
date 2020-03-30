package cmd

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/kondoumh/scrapbox-viz/pkg/file"

	"github.com/kondoumh/scrapbox-viz/pkg/api"
	"github.com/kondoumh/scrapbox-viz/pkg/types"
	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch all pages of the project",
	Long: LongUsage(`
		Fetch all page data of the project.

		  sbv fetch -p <project name>

		Page list data will be saved as JSON file at '<WorkDir>/<project name>.json'.
		Each Page data will be saved as JSON file in '<WorkDir>/<project name>'.
		The file name consists of the page ID.
	`),
	Run: func(cmd *cobra.Command, args []string) {
		doFetch(cmd)
	},
}

func init() {
	fetchCmd.PersistentFlags().StringP("project", "p", "help-jp", "Name of Scrapbox project")
	rootCmd.AddCommand(fetchCmd)
}

func doFetch(cmd *cobra.Command) {
	projectName, _ := cmd.PersistentFlags().GetString("project")
	project, err := fetchIndex(projectName)
	CheckErr(err)
	fmt.Printf("fetch all pages, %s : %d\n", project.Name, project.Count)
	err = fetchPageList(project)
	CheckErr(err)
	groups, err := dividePagesList(3, projectName)
	CheckErr(err)
	path := fmt.Sprintf("%s/%s", config.WorkDir, projectName)
	file.CreateDir(path)
	var wg sync.WaitGroup
	start := time.Now()
	wg.Add(len(groups))
	for _, pages := range groups {
		go fetchPagesByGroup(projectName, pages, &wg)
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("took %s\n", elapsed)
}

func fetchIndex(projectName string) (types.Project, error) {
	data, err := api.FetchIndex(projectName)
	var project types.Project
	if err != nil {
		return project, err
	}
	err = json.Unmarshal(data, &project)
	if err != nil {
		return project, err
	}
	return project, nil
}

func fetchPageList(project types.Project) error {
	pages := []types.Page{}
	for skip := 0; skip < project.Count; skip += api.Limit {
		data, err := api.FetchPageList(project.Name, skip)
		if err != nil {
			return err
		}
		var proj types.Project
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

func dividePagesList(multiplicity int, projectName string) ([][]types.Page, error) {
	var divided [][]types.Page
	var proj types.Project
	err := proj.ReadFrom(projectName, config.WorkDir)
	if err != nil {
		return divided, err
	}
	fmt.Printf("Total pages : %d\n", len(proj.Pages))
	chunkSize := len(proj.Pages) / multiplicity
	fmt.Printf("Chunk size : %d\n", chunkSize)
	for i := 0; i < len(proj.Pages); i += chunkSize {
		end := i + chunkSize
		if end > len(proj.Pages) {
			end = len(proj.Pages)
		}
		divided = append(divided, proj.Pages[i:end])
	}
	totalCount := 0
	for _, pages := range divided {
		totalCount += len(pages)
		fmt.Printf("Size of chunk %d\n", len(pages))
	}
	fmt.Printf("Total pages to be fetched %d\n", totalCount)
	return divided, nil
}

func fetchPagesByGroup(projectName string, pages []types.Page, wg *sync.WaitGroup) error {
	defer wg.Done()
	for _, page := range pages {
		fmt.Println(page.Title)
		if err := fetchPage(projectName, page.Title, page.ID); err != nil {
			return err
		}
	}
	return nil
}

func fetchPage(projectName string, title string, index string) error {
	data, err := api.FetchPage(projectName, title)
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("%s.json", index)
	path := fmt.Sprintf("%s/%s", config.WorkDir, projectName)
	if err := file.WriteBytes(data, fileName, path); err != nil {
		return err
	}
	return nil
}
