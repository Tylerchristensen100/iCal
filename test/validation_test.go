package test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	ical "github.com/Tylerchristensen100/iCal"
	"github.com/Tylerchristensen100/iCal/timezones"
)

func TestCalendarOutput(t *testing.T) {
	cal := ical.Create("Test Calendar", "This is a test calendar.")
	err := cal.AddEvent(ical.Event{Title: "Test Event",
		Description: "This is a test event.",
		StartDate:   time.Date(2024, 8, 1, 10, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, 8, 1, 11, 0, 0, 0, time.UTC),
		TimeZone:    ical.TimeZone(timezones.US_Mountain),
	})
	if err != nil {
		t.Fatalf("Failed to add event: %v", err)
	}

	organizer := ical.Participant{
		Name:  "Test Person",
		Email: "test@test.org",
	}

	dueDate := time.Date(2026, 1, 2, 12, 0, 0, 0, time.UTC)
	err = cal.AddTodo(ical.Todo{
		Summary:     "Test Todo",
		Description: "This is a test todo item.",
		Due:         &dueDate,
		Status:      ical.InProcessStatus,
		Organizer:   organizer,
	})
	if err != nil {
		t.Fatalf("Failed to add todo: %v", err)
	}

	err = cal.AddJournal(ical.Journal{
		Summary:     "Test Journal",
		Description: "This is a test journal entry.",
		Status:      ical.DraftJournal,
		Organizer:   organizer,
	})
	if err != nil {
		t.Fatalf("Failed to add journal: %v", err)
	}

	data, err := cal.Generate()
	if err != nil {
		t.Fatalf("Failed to generate calendar: %v", err)
	}
	err = cal.Save(fmt.Sprintf("./tmp/test_calendar_output%s.ics", time.Now().Format("20060102150405")))
	if err != nil {
		t.Fatalf("Failed to save calendar: %v", err)
	}

	if len(data) == 0 {
		t.Errorf("Generated calendar data is empty")
	}

	expectedStart := "BEGIN:VCALENDAR"
	if !strings.HasPrefix(string(data), expectedStart) {
		t.Errorf("Expected calendar to start with '%s', got '%s'", expectedStart, data[:len(expectedStart)])
	}

	expectedEnd := "END:VCALENDAR\r\n"
	calSuffix := string(data[len(data)-len(expectedEnd):])
	if calSuffix != expectedEnd {
		t.Errorf("Expected calendar to end with '%s', got '%s'", expectedEnd, calSuffix)
	}

}
