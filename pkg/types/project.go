package types

import (
	"encoding/json"

	"github.com/mamezou-tech/sbgraph/pkg/file"
)

// Project represents a Scrapbox project
type Project struct {
	Name  string `json:"projectName"`
	Date  string `json:"date"`
	Count int    `json:"count"`
	Skip  int    `json:"skip"`
	Pages []Page `json:"pages"`
}

// ReadFrom will deserialize Project from file
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
