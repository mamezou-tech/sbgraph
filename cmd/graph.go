package cmd

import (
	"fmt"
	"os"

	"github.com/kondoumh/sbgraph/pkg/types"
	"github.com/mzohreva/GoGraphviz/graphviz"
	"github.com/spf13/cobra"
)

// graphCmd represents the graph command
var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Generate graph structure of the project",
	Long: LongUsage(`
		Generate graph structure of pages and authors.

		  sbf graph

		Graphviz dot file will be created at '<WorkDir>/<project name>.dot'.
	`),
	Run: func(cmd *cobra.Command, args []string) {
		buildGraph(cmd)
	},
}

func init() {
	graphCmd.PersistentFlags().IntP("threshold", "t", 0, "Threshold value of views to filter page")
	graphCmd.PersistentFlags().BoolP("include", "i", false, "Include user node")
	graphCmd.PersistentFlags().BoolP("anonymize", "a", false, "Anonymize user")

	rootCmd.AddCommand(graphCmd)
}

func buildGraph(cmd *cobra.Command) {
	projectName := config.CurrentProject
	CheckProject(projectName)
	threshold, _ := cmd.PersistentFlags().GetInt("threshold")
	includeUser, _ := cmd.PersistentFlags().GetBool("include")
	anonymize, _ := cmd.PersistentFlags().GetBool("anonymize")
	fmt.Printf("Build graph project : %s, threshold : %d, include user: %t, anonymize : %t\n", projectName, threshold, includeUser, anonymize)
	var proj types.Project
	err := proj.ReadFrom(projectName, config.WorkDir)
	CheckErr(err)

	users := map[string]types.User{}
	pages := map[string]types.Page{}
	for _, idx := range proj.Pages {
		if idx.Views <= threshold {
			continue
		}
		var page types.Page
		err := page.ReadFrom(projectName, idx.ID, config.WorkDir)
		CheckErr(err)
		pages[page.ID] = page
		u, contains := users[page.Author.ID]
		if contains {
			u.PagesCreated = append(u.PagesCreated, page.ID)
			users[page.Author.ID] = u
		} else {
			user := types.User{
				ID:           page.Author.ID,
				Name:         page.Author.Name,
				DisplayName:  page.Author.DisplayName,
				PagesCreated: []string{page.ID},
			}
			users[page.Author.ID] = user
		}
	}

	graph := createGraph()
	pNodes := map[string]int{}
	for _, p := range pages {
		gid := graph.AddNode(p.Title)
		pNodes[p.ID] = gid
	}
	uNodes := map[string]int{}
	if includeUser {
		for _, u := range users {
			if len(u.PagesCreated) == 0 {
				continue
			}
			var username string
			if anonymize {
				username = u.ID[5:10]
			} else {
				username = u.DisplayName
			}
			gid := graph.AddNode(username)
			graph.NodeAttribute(gid, graphviz.FillColor, "cyan")
			uNodes[u.ID] = gid
		}
	}
	for _, p := range pages {
		pid, _ := pNodes[p.ID]
		for _, link := range p.Related.Links {
			lid, contains := pNodes[link.ID]
			if contains {
				graph.AddEdge(pid, lid, "")
			}
		}
	}
	if includeUser {
		for _, u := range users {
			uid, _ := uNodes[u.ID]
			for _, c := range u.PagesCreated {
				pid, contains := pNodes[c]
				if contains {
					graph.AddEdge(uid, pid, "")
				}
			}
		}
	}
	err = writeDot(graph, projectName, config.WorkDir)
	CheckErr(err)
}

func createGraph() graphviz.Graph {
	graph := graphviz.Graph{}
	graph.MakeDirected()

	graph.GraphAttribute(graphviz.NodeSep, "0.5")

	graph.DefaultNodeAttribute(graphviz.Shape, graphviz.ShapeBox)
	graph.DefaultNodeAttribute(graphviz.FontName, "Courier")
	graph.DefaultNodeAttribute(graphviz.FontSize, "14")
	graph.DefaultNodeAttribute(graphviz.Style, graphviz.StyleFilled+","+graphviz.StyleRounded)
	graph.DefaultNodeAttribute(graphviz.FillColor, "yellow")
	graph.DefaultEdgeAttribute(graphviz.FontName, "Courier")
	graph.DefaultEdgeAttribute(graphviz.FontSize, "12")

	return graph
}

func writeDot(graph graphviz.Graph, projectName string, workDir string) error {
	path := fmt.Sprintf("%s/%s.dot", workDir, projectName)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	graph.GenerateDOT(file)
	return nil
}
