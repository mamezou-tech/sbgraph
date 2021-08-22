package cmd

import "testing"

func TestContains(t *testing.T) {
	strs := []string{"Foo", "Bar", "あああ"}
	got := Contains(strs, "foo")
	if !got {
		t.Errorf("Contains() %v, want %v", got, true)
	}
	got = Contains(strs, "BAR")
	if !got {
		t.Errorf("Contains() %v, want %v", got, true)
	}
	got = Contains(strs, "あああ")
	if !got {
		t.Errorf("Contains() %v, want %v", got, true)
	}
}
