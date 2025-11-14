package ical

import "testing"

func TestCleanDescription(t *testing.T) {
	rawDescription := "This is a test description with especially long content, and a lot of words and special characters."
	cleanedDescription := cleanDescription(rawDescription)
	expectedDescription := "This is a test description with especially long content, and a \r\n lot of words and special characters."
	if cleanedDescription != expectedDescription {
		t.Errorf("Expected cleaned description to be:\n%s\nGot:\n%s", expectedDescription, cleanedDescription)
	}

	clean := "Short description."
	cleaned := cleanDescription(clean)
	if cleaned != clean {
		t.Errorf("Expected cleaned description to be unchanged:\n%s\nGot:\n%s", clean, cleaned)
	}

}

func TestEscapeText(t *testing.T) {
	dirtyText := "This is a test; with, special\ncharacters: \\ and more."
	expected := "This is a test\\; with\\, special\\ncharacters: \\\\ and more."
	escaped := escapeText(dirtyText)
	if escaped != expected {
		t.Errorf("escapeText() = %v, want %v", escaped, expected)
	}
}
