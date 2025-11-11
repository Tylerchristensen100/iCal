package ical

import "testing"

func TestCreate(t *testing.T) {
	const (
		name        = "Test Calendar"
		description = "This is a test calendar."
	)
	cal := Create(name, description)
	if cal.Name != name {
		t.Errorf("Expected calendar name to be '%s', got '%s'", name, cal.Name)
	}
	if cal.Description != description {
		t.Errorf("Expected calendar description to be '%s', got '%s'", description, cal.Description)
	}
	if len(cal.Events) != 0 {
		t.Errorf("Expected calendar to have 0 events, got %d", len(cal.Events))
	}

	calFail := Create("", "")
	if calFail != nil {
		t.Errorf("Expected nil for empty calendar name and description, got '%v'", calFail)
	}

}
