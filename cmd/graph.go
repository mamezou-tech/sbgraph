package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mamezou-tech/sbgraph/pkg/file"

	"github.com/mamezou-tech/sbgraph/pkg/types"
	"github.com/mzohreva/GoGraphviz/graphviz"
	"github.com/spf13/cobra"
)

type page struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Views  int    `json:"views"`
	Linked int    `json:"linked"`
}

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type pageLink struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type userPage struct {
	UserID string `json:"user"`
	PageID string `json:"page"`
}

type projectGraph struct {
	Pages     []page     `json:"pages"`
	Users     []user     `json:"users"`
	Links     []pageLink `json:"links"`
	UserPages []userPage `json:"userPages"`
}

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
	graphCmd.PersistentFlags().BoolP("json", "j", false, "Output as JSON format")
	graphCmd.PersistentFlags().BoolP("image", "m", false, "Output SVG image")

	rootCmd.AddCommand(graphCmd)
}

func buildGraph(cmd *cobra.Command) {
	projectName := config.CurrentProject
	CheckProject(projectName)
	threshold, _ := cmd.PersistentFlags().GetInt("threshold")
	includeUser, _ := cmd.PersistentFlags().GetBool("include")
	anonymize, _ := cmd.PersistentFlags().GetBool("anonymize")
	oJSON, _ := cmd.PersistentFlags().GetBool("json")
	oSVG, _ := cmd.PersistentFlags().GetBool("image")

	fmt.Printf("Build graph project : %s, threshold : %d, include user: %t, anonymize : %t, json : %t, svg : %t\n", projectName, threshold, includeUser, anonymize, oJSON, oSVG)
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
	var pGraph projectGraph
	pNodes := map[string]int{}
	for _, p := range pages {
		gid := graph.AddNode(p.Title)
		pNodes[p.ID] = gid
		pGraph.Pages = append(pGraph.Pages, page{p.ID, p.Title, p.Views, p.Linked})
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
			pGraph.Users = append(pGraph.Users, user{u.ID, username})
		}
	}
	for _, p := range pages {
		pid, _ := pNodes[p.ID]
		for _, link := range p.Related.Links {
			lid, contains := pNodes[link.ID]
			if contains {
				graph.AddEdge(pid, lid, "")
				pGraph.Links = append(pGraph.Links, pageLink{p.ID, link.ID})
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
					pGraph.UserPages = append(pGraph.UserPages, userPage{u.ID, c})
				}
			}
		}
	}
	err = writeDot(graph, projectName, config.WorkDir)
	CheckErr(err)
	if oJSON {
		data, _ := json.Marshal(pGraph)
		err = file.WriteBytes(data, projectName+"_graph.json", config.WorkDir)
		CheckErr(err)
	}
	if oSVG {
		err = writeSvg(graph, projectName, config.WorkDir)
		CheckErr(err)
	}
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

func writeSvg(graph graphviz.Graph, projectName string, workDir string) error {
	path := fmt.Sprintf("%s/%s.svg", workDir, projectName)
	return graph.GenerateImage("dot", path, "svg")
}
