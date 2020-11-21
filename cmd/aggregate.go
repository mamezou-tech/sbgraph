package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mamezou-tech/sbgraph/pkg/file"

	"github.com/cheggaaa/pb/v3"
	"github.com/mamezou-tech/sbgraph/pkg/types"
	"github.com/spf13/cobra"
)

// aggregateCmd represents the aggregate command
var aggregateCmd = &cobra.Command{
	Use:   "aggregate",
	Short: "Aggregate project activities",
	Long: LongUsage(`
	Aggregate project activities.

	  sbf aggregate

	JSON will be created as '<WorkDir>/<project name>_ag.json'.
	If the csv flag is specified, CSV will be created as '<WorkDir>/project name>_ag.csv'.
	`),
	Run: func(cmd *cobra.Command, args []string) {
		doAggregate(cmd)
	},
}

func init() {
	aggregateCmd.PersistentFlags().BoolP("csv", "s", false, "Output as CSV")
	rootCmd.AddCommand(aggregateCmd)
}

func doAggregate(cmd *cobra.Command) {
	csv, _ := cmd.PersistentFlags().GetBool("csv")
	projectName := config.CurrentProject
	CheckProject(projectName)
	fmt.Printf("Aggregate project : %s\n", projectName)
	var proj types.Project
	err := proj.ReadFrom(projectName, config.WorkDir)
	CheckErr(err)
	contrib := map[string]types.Contribution{}
	bar := pb.StartNew(proj.Count)
	for _, idx := range proj.Pages {
		var page types.Page
		err := page.ReadFrom(projectName, idx.ID, config.WorkDir)
		CheckErr(err)
		p, contains := contrib[page.Author.ID]
		if contains {
			p.PagesCreated++
			p.ViewsCreatedPages += page.Views
			p.LinksCreatedPages += page.Linked
			contrib[page.Author.ID] = p
		} else {
			c := types.Contribution{
				UserID:            page.Author.ID,
				UserName:          page.Author.DisplayName,
				PagesContributed:  1,
				ViewsCreatedPages: page.Views,
				LinksCreatedPages: page.Linked,
			}
			contrib[page.Author.ID] = c
		}
		for _, user := range page.Collaborators {
			p, contains := contrib[user.ID]
			if contains {
				p.PagesContributed++
				contrib[user.ID] = p
			} else {
				c := types.Contribution{
					UserID:           user.ID,
					UserName:         user.DisplayName,
					PagesContributed: 1,
				}
				contrib[user.ID] = c
			}
		}
		bar.Increment()
	}
	bar.Finish()
	err = writeContrib(projectName, contrib, csv)
	CheckErr(err)
}

func writeContrib(projectName string, contrib map[string]types.Contribution, csv bool) error {
	if csv {
		path := fmt.Sprintf("%s/%s_contrib.csv", config.WorkDir, projectName)
		fmt.Println(path)
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
		file.Write(([]byte)("User ID, User Name,Pages Created,Pages Contributed,Views of Created Pages,Links of Created Pages\n"))
		for _, v := range contrib {
			data := fmt.Sprintf("%s,%s,%d,%d,%d,%d\n", v.UserID, v.UserName, v.PagesCreated, v.PagesContributed, v.ViewsCreatedPages, v.LinksCreatedPages)
			_, err = file.Write(([]byte)(data))
			if err != nil {
				return err
			}
		}
	} else {
		data, _ := json.Marshal(contrib)
		if err := file.WriteBytes(data, projectName+"_contrib.json", config.WorkDir); err != nil {
			return err
		}
	}
	return nil
}
