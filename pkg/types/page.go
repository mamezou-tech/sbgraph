package types

import (
	"encoding/json"

	"github.com/kondoumh/scrapbox-viz/pkg/file"
)

// Page represents a Scrapbox page
type Page struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Views         int    `json:"views"`
	Linked        int    `json:"linked"`
	Author        User   `json:"user"`
	Collaborators []User `json:"collaborators"`
	Related       struct {
		Links []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"links1hop"`
	} `json:"relatedPages"`
}

// User represents a Scrapbox user
type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
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
