package cmd

import (
	"fmt"
	"os"

	"github.com/cheggaaa/pb/v3"
	"github.com/kondoumh/sbgraph/pkg/types"
	"github.com/spf13/cobra"
)

// aggregateCmd represents the aggregate command
var aggregateCmd = &cobra.Command{
	Use:   "aggregate",
	Short: "Aggregate project activities",
	Long: LongUsage(`
	Aggregate project activities.

	  sbf aggregate -p <project name>

	CSV will be created at '<WorkDir>/<project name>.csv'.
	`),
	Run: func(cmd *cobra.Command, args []string) {
		doAggregate(cmd)
	},
}

func init() {
	aggregateCmd.PersistentFlags().StringP("project", "p", "help-jp", "Name of Scrapbox project")
	rootCmd.AddCommand(aggregateCmd)
}

type contribute struct {
	UserID            string
	UserName          string
	PagesCreated      int
	PagesContributed  int
	ViewsCreatedPages int
	LinksCreatedPages int
}

func doAggregate(cmd *cobra.Command) {
	projectName, _ := cmd.PersistentFlags().GetString("project")
	fmt.Printf("Aggregate project : %s\n", projectName)
	var proj types.Project
	err := proj.ReadFrom(projectName, config.WorkDir)
	CheckErr(err)
	contrib := map[string]contribute{}
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
			c := contribute{
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
				c := contribute{
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
	err = writeContrib(projectName, contrib)
	CheckErr(err)
}

func writeContrib(projectName string, contrib map[string]contribute) error {
	path := fmt.Sprintf("%s/%s.csv", config.WorkDir, projectName)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(([]byte)("User Name,Pages Created,Pages Contributed,Views of Created Pages,Links of Created Pages\n"))
	for _, v := range contrib {
		data := fmt.Sprintf("%s,%d,%d,%d,%d\n", v.UserName, v.PagesCreated, v.PagesContributed, v.ViewsCreatedPages, v.LinksCreatedPages)
		_, err = file.Write(([]byte)(data))
		if err != nil {
			return err
		}
	}
	return nil
}
