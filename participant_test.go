package ical

import (
	"fmt"
	"strings"
	"testing"
)

func TestEventGeneration(t *testing.T) {
	var tests = []struct {
		name        string
		email       string
		expectError bool
	}{
		{"Valid User", "test@example.com", false},
		{"Invalid User", "invalid-email", true},
	}
	for _, tt := range tests {
		event := mockEvent()
		event.AddAttendee(tt.name, tt.email)
		builder := strings.Builder{}
		err := event.buildEventDetails(&builder)
		if err != nil {
			t.Fatalf("buildEventDetails() returned error: %v", err)
		}
		result := builder.String()
		if !strings.Contains(result, fmt.Sprintf("CN=%s;", tt.name)) && !strings.Contains(result, fmt.Sprintf("mailto:%s", tt.email)) && !tt.expectError {
			t.Errorf("Expected attendee %s with email %s to be in event details, but it was not found", tt.name, tt.email)
		}
	}
}

func TestGenerate(t *testing.T) {
	const (
		name  = "Attendee Name"
		email = "attendee@example.com"
	)
	var builder strings.Builder
	p := Participant{Name: name, Email: email}
	err := p.generateOrganizer(&builder)
	if err != nil {
		t.Fatalf("AddOrganizer() returned error: %v", err)
	}

	result := builder.String()
	if !strings.Contains(result, fmt.Sprintf("CN=%s;", name)) && !strings.Contains(result, fmt.Sprintf("mailto:%s", email)) {
		t.Errorf("Expected attendee %s with email %s to be in event details, but it was not found", name, email)
	}
}

func TestGenerateOrganizer(t *testing.T) {
	const (
		name  = "Organizer Name"
		email = "organizer@example.com"
	)
	var builder strings.Builder
	p := Participant{Name: name, Email: email}
	err := p.generateOrganizer(&builder)
	if err != nil {
		t.Fatalf("AddOrganizer() returned error: %v", err)
	}

	result := builder.String()
	if !strings.Contains(result, "CN=Organizer Name:mailto:organizer@example.com") {
		t.Errorf("Expected organizer %s with email %s to be in event details, but it was not found", name, email)
	}
}
