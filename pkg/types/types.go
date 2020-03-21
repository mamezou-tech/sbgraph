package types

import (
	"encoding/json"

	"github.com/kondoumh/scrapbox-viz/pkg/file"
)

// Project represents a Scrapbox project
type Project struct {
	Name  string `json:"projectName"`
	Count int    `json:"count"`
	Skip  int    `json:"skip"`
	Pages []Page `json:"pages"`
}

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

// ReadFrom will serialize Project from file
func (project *Project) ReadFrom(projectName string, workDir string) error {
	bytes, err := file.ReadBytes(projectName+".json", workDir)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, &project); err != nil {
		return err
	}
	return nil
}
