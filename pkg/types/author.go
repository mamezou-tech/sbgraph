package types

import (
	"encoding/json"

	"github.com/mamezou-tech/sbgraph/pkg/file"
)

// Authors represents authors of the Scrapbox project
type Authors struct {
	Authors []Author `json:"authors"`
}

// Author represents author of the Scrapbox page
type Author struct {
	ID string `json:"id"`
	Name string `json:"name"`
}

// ReadFrom will deserialize Authors from file
func (authors *Authors) ReadFrom(projectName string, workDir string) error {
	bytes, err := file.ReadBytes(projectName+"_authors.json", workDir)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, &authors); err != nil {
		return err
	}
	return nil
}
