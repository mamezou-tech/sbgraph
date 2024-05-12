package cmd

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"github.com/mamezou-tech/sbgraph/pkg/file"
	"github.com/mamezou-tech/sbgraph/pkg/types"
	"github.com/spf13/cobra"
)

type pageSimple struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Authors []author `json:"authors"`
	Lines   []string `json:"lines"`
}

type author struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// extractCmd represents the extract command
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract from downloaded JSON files",
	Long: LongUsage(`Extract from downloaded JSON files that matches passed tag.

  	sbgraph extract -i "foo bar baz" -e "hoge huga"`),
	Run: func(cmd *cobra.Command, args []string) {
		doExtract(cmd)
	},
}

func init() {
	extractCmd.PersistentFlags().StringP("includes", "i", "", "Words for extracting pages(space delimited).")
	extractCmd.PersistentFlags().StringP("excludes", "e", "", "Words to exclude when extracting pages(space delimited).")
	extractCmd.PersistentFlags().StringP("suffix", "s", "extracted", "suffix for output directory")
	rootCmd.AddCommand(extractCmd)
}

func doExtract(cmd *cobra.Command) {
	projectName := config.CurrentProject
	CheckProject(projectName)
	tagsStr, _ := cmd.PersistentFlags().GetString("tags")
	excludesStr, _ := cmd.PersistentFlags().GetString("excludes")
	suffix, _ := cmd.PersistentFlags().GetString("suffix")

	includes := strings.Split(tagsStr, " ")
	excludes := strings.Split(excludesStr, " ")
	outputDir := projectName + "-" + suffix

	fmt.Printf("Extract files : %s, tags : %s, excludes : %s, output: %s\n", projectName, includes, excludes, outputDir)
	var proj types.Project
	err := proj.ReadFrom(projectName, config.WorkDir)
	CheckErr(err)

	bar := pb.StartNew(proj.Count)

	var authors types.Authors
	if file.Exists(projectName+"_authors.json", config.WorkDir) {
		err := authors.ReadFrom(projectName, config.WorkDir)
		CheckErr(err)
	}

	outputPath := config.WorkDir + "/" + outputDir
	file.CreateDir(outputPath)
	for _, idx := range proj.Pages {
		var page types.Page
		err := page.ReadFrom(projectName, idx.ID, config.WorkDir)
		CheckErr(err)
		result := isExtractable(toLines(&page), includes, excludes)
		if result {
			var simplePage pageSimple
			simplePage.ID = page.ID
			simplePage.Title = page.Title
			simplePage.Authors = []author{}
			displayName := getDisplayName(authors, page.Author.ID, page.Author.Name)
			simplePage.Authors = append(simplePage.Authors, author{ID: page.Author.ID, Name: displayName})
			for _, collaborator := range page.Collaborators {
				displayName := getDisplayName(authors, collaborator.ID, collaborator.Name)
				simplePage.Authors = append(simplePage.Authors, author{ID: collaborator.ID, Name: displayName})
			}
			simplePage.Lines = toLines(&page)
			data, _ := json.Marshal(simplePage)
			err = file.WriteBytes(data, simplePage.ID+".json", outputPath)
			CheckErr(err)
		}
		bar.Increment()
	}
	bar.Finish()
}

func toLines(page *types.Page) []string {
	lines := []string{}
	for _, line := range page.Lines {
		lines = append(lines, line.Text)
	}
	return lines
}

func isExtractable(lines []string, includes []string, excludes []string) bool {
	if !isEmpty(excludes) {
		for _, exclude := range excludes {
			re := regexp.MustCompile(exclude)
			for _, line := range lines {
				if re.MatchString(line) {
					return false
				}
			}
		}
	}
	if isEmpty(includes) {
		return true
	}
	for _, include := range includes {
		re := regexp.MustCompile(include)
		for _, line := range lines {
			if re.MatchString(line) {
				return true
			}
		}
	}
	return false
}

func isEmpty(arr []string) bool {
	return len(arr) == 0 || (len(arr) == 1 && arr[0] == "")
}

func getDisplayName(authors types.Authors, id string, name string) string {
	username := name
	for _, a := range authors.Authors {
		if id == a.ID {
			username = a.Name
			break
		}
	}
	return username
}
