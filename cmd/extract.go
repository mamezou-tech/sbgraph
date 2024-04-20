package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"github.com/mamezou-tech/sbgraph/pkg/file"
	"github.com/mamezou-tech/sbgraph/pkg/types"
	"github.com/spf13/cobra"
)

type pageSimple struct {
	ID		string `json:"id"`
	Title	string `json:"title"`
	Lines	[]string `json:"lines"`
}

// extractCmd represents the extract command
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract from downloaded JSON files",
	Long: LongUsage(`Extract from downloaded JSON files that matches passed tag.

  	sbgraph extract -t tagname -o outputdir`),
	Run: func(cmd *cobra.Command, args []string) {
		doExtract(cmd);
	},
}

func init() {
	extractCmd.PersistentFlags().StringP("tag", "t", "", "Extract pages with the specified tag.")
	extractCmd.PersistentFlags().StringP("suffix", "s", "", "suffix for output directory")
	rootCmd.AddCommand(extractCmd)
}

func doExtract(cmd *cobra.Command) {
	projectName := config.CurrentProject
	CheckProject(projectName)
	tag, _ := cmd.PersistentFlags().GetString("tag")
	suffix, _ := cmd.PersistentFlags().GetString("suffix")
	CheckArg(tag, "tag");
	CheckArg(suffix, "suffix");

	fmt.Printf("Extract files : %s, tag : %s\n", projectName, tag)
	var proj types.Project
	err := proj.ReadFrom(projectName, config.WorkDir)
	CheckErr(err)

	bar := pb.StartNew(proj.Count)

	outputDir := config.WorkDir + "/" + projectName + "-" + suffix
	file.CreateDir(outputDir)
	for _, idx := range proj.Pages {
		var page types.Page
		err := page.ReadFrom(projectName, idx.ID, config.WorkDir)
		CheckErr(err)
		result := containsTag(toLines(&page), tag)
		if result {
			var simplePage pageSimple
			simplePage.ID = page.ID
			simplePage.Title = page.Title
			simplePage.Lines = toLines(&page)
			data, _ := json.Marshal(simplePage)
			err = file.WriteBytes(data, simplePage.ID+".json", outputDir)
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

func containsTag(lines []string, tag string) bool {
	for _, line := range lines {
		if strings.Contains(line, tag) {
			return true
		}
	}
	return false
}
