package types

import (
	"encoding/json"

	"github.com/mamezou-tech/sbgraph/pkg/file"
)

// Page represents a Scrapbox page
type Page struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Created       int32    `json:"created"`
	Updated       int32    `json:"updated"`
	Pin           int64    `json:"pin"`
	Views         int      `json:"views"`
	Linked        int      `json:"linked"`
	Author        User     `json:"user"`
	Collaborators []User   `json:"collaborators"`
	Image         string   `json:"image"`
	Tags          []string `json:"links"`
	Related       struct {
		Links []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"links1hop"`
	} `json:"relatedPages"`
}

// User represents a Scrapbox user
type User struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	DisplayName  string `json:"displayName"`
	PagesCreated []string
}

// Contribution represents summary of the user's contribution
type Contribution struct {
	UserID            string `json:"userId"`
	UserName          string `json:"userName"`
	PagesCreated      int    `json:"pagesCreated"`
	PagesContributed  int    `json:"pagesContributed"`
	ViewsCreatedPages int    `json:"viewsCreatedPages"`
	LinksCreatedPages int    `json:"linksCreatedPages"`
}

// ReadFrom will deserialize Project from file
func (page *Page) ReadFrom(projectName string, id string, workDir string) error {
	bytes, err := file.ReadBytes(id+".json", workDir+"/"+projectName)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, &page); err != nil {
		return err
	}
	return nil
}
