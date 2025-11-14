package test

import (
	"strings"
	"testing"
	"time"

	ical "github.com/Tylerchristensen100/iCal"
	"github.com/Tylerchristensen100/iCal/timezones"
)

// Test TODO, JOURNAL, VTIMEZONE, ALARM, EVENT components
func TestValidComponents(t *testing.T) {
	cal := ical.Create("Test Calendar", "This is a test calendar.")

	journal := ical.Journal{
		Summary:     "Test Journal",
		Description: "This is a test journal entry.",
	}
	err := cal.AddJournal(journal)
	if err != nil {
		t.Fatalf("Failed to add journal: %v", err)
	}

	todo := ical.Todo{
		Summary:     "Test Todo",
		Description: "This is a test todo item.",
	}
	err = cal.AddTodo(todo)
	if err != nil {
		t.Fatalf("Failed to add todo: %v", err)
	}

	event := ical.Event{
		Title:       "Test Event",
		Description: "This is a test event.",
		StartDate:   time.Date(2024, 8, 1, 10, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, 8, 1, 11, 0, 0, 0, time.UTC),
		TimeZone:    ical.TimeZone(timezones.US_Mountain),
	}

	event.AddReminder(
		ical.Reminder{
			Action:      ical.DisplayReminderAction,
			Description: "Test Reminder",
			Trigger:     time.Duration(-15 * time.Minute),
		},
	)
	err = cal.AddEvent(event)
	if err != nil {
		t.Fatalf("Failed to add event: %v", err)
	}

	data, err := cal.Generate()
	if err != nil {
		t.Fatalf("Failed to generate calendar with components: %v", err)
	}

	if len(data) == 0 {
		t.Errorf("Generated calendar data is empty")
	}

	if !strings.Contains(string(data), "BEGIN:VJOURNAL") {
		t.Errorf("Generated calendar missing VJOURNAL component")
	}

	if !strings.Contains(string(data), "BEGIN:VTODO") {
		t.Errorf("Generated calendar missing VTODO component")
	}

	if !strings.Contains(string(data), "BEGIN:VTIMEZONE") {
		t.Errorf("Generated calendar missing VTIMEZONE component")
	}

	if !strings.Contains(string(data), "BEGIN:VALARM") {
		t.Errorf("Generated calendar missing VALARM component")
	}

	if !strings.Contains(string(data), "BEGIN:VEVENT") {
		t.Errorf("Generated calendar missing VEVENT component")
	}

}

func TestInvalidComponents(t *testing.T) {
	cal := ical.Create("Invalid Calendar", "This calendar has invalid components.")

	invalidJournal := ical.Journal{
		Summary: "",
	}
	err := cal.AddJournal(invalidJournal)
	if err == nil {
		t.Errorf("Expected error when adding invalid journal, got nil")
	}

	invalidTodo := ical.Todo{
		Summary: "",
	}
	err = cal.AddTodo(invalidTodo)
	if err == nil {
		t.Errorf("Expected error when adding invalid todo, got nil")
	}

	invalidEvent := ical.Event{
		Title:     "Invalid Event",
		StartDate: time.Date(2024, 8, 1, 12, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 5, 1, 11, 0, 0, 0, time.UTC), // End before start

	}
	err = cal.AddEvent(invalidEvent)
	if err == nil {
		t.Errorf("Expected error when adding invalid event, got nil")
	}

	data, err := cal.Generate()
	if err == nil {
		t.Errorf("Expected error when generating calendar with invalid components, got nil")
	}

	if len(data) != 0 {
		t.Errorf("Expected empty calendar data when generation fails, got data of length %d", len(data))
	}
}
