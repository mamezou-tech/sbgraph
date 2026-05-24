package types

import (
	"encoding/json"
	"testing"
)

func TestPageUnmarshalJSONSupportsLegacyCollaborators(t *testing.T) {
	data := []byte(`{"id":"page-id","title":"Page","user":{"id":"author-id"},"collaborators":[{"id":"collab-id","name":"collab-name"}]}`)

	var page Page
	if err := json.Unmarshal(data, &page); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if len(page.Collaborators) != 1 {
		t.Fatalf("len(page.Collaborators) = %d, want 1", len(page.Collaborators))
	}
	if page.Collaborators[0].ID != "collab-id" {
		t.Fatalf("page.Collaborators[0].ID = %q, want %q", page.Collaborators[0].ID, "collab-id")
	}
}
