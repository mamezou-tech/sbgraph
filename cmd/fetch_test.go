package cmd

import (
	"testing"

	"github.com/mamezou-tech/sbgraph/pkg/types"
)

func TestEnrichPageUsers(t *testing.T) {
	page := types.Page{
		Author: types.User{ID: "author-id"},
		Collaborators: []types.User{
			{ID: "collab-id"},
		},
	}
	projectUsers := map[string]types.User{
		"author-id": {
			ID:          "author-id",
			Name:        "author-name",
			DisplayName: "author-display",
		},
		"collab-id": {
			ID:          "collab-id",
			Name:        "collab-name",
			DisplayName: "collab-display",
		},
	}

	enrichPageUsers(&page, projectUsers)

	if page.Author.Name != "author-name" || page.Author.DisplayName != "author-display" {
		t.Fatalf("author was not enriched: %#v", page.Author)
	}
	if len(page.Collaborators) != 1 {
		t.Fatalf("unexpected collaborators: %#v", page.Collaborators)
	}
	if page.Collaborators[0].Name != "collab-name" || page.Collaborators[0].DisplayName != "collab-display" {
		t.Fatalf("collaborator was not enriched: %#v", page.Collaborators[0])
	}
}
