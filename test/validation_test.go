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

func TestCalendarWithMultipleOccurrences(t *testing.T) {
	cal := ical.Create("Multiple Occurrences Calendar", "Calendar with multiple occurrences event.")
	recurrence := ical.Recurrences{
		Frequency: ical.WeeklyFrequency,
		Day:       time.Wednesday,
		StartTime: time.Date(0, 0, 0, 14, 0, 0, 0, time.UTC),
		EndTime:   time.Date(0, 0, 0, 15, 0, 0, 0, time.UTC),
	}

	err := cal.AddEvent(ical.Event{
		Title:       "Weekly Meeting",
		Description: "This is a weekly meeting.",
		Recurrences: []ical.Recurrences{recurrence},
		TimeZone:    ical.TimeZone(timezones.US_Mountain),
		StartDate:   time.Date(2025, 9, 3, 14, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2026, 12, 3, 15, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Failed to add event with recurrences: %v", err)
	}

	err = cal.AddEvent(ical.Event{
		Title:       "Single Event",
		Description: "This is a single event.",
		TimeZone:    ical.TimeZone(timezones.US_Mountain),
		StartDate:   time.Now().Add(24 * time.Hour),
		EndDate:     time.Now().Add(25 * time.Hour),
		URL:         "https://example.com/single-event",
	})
	if err != nil {
		t.Fatalf("Failed to add single event: %v", err)
	}

	err = cal.AddEvent(ical.Event{
		Title:       "Multiple Times a Week Event",
		Description: "This Event Happens Tuesday, Thursday, and Every other Friday",
		StartDate:   time.Date(2025, 8, 1, 9, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2025, 12, 20, 10, 0, 0, 0, time.UTC),
		Recurrences: []ical.Recurrences{
			//Tuesday
			{
				Frequency: ical.WeeklyFrequency,
				Day:       time.Tuesday,
				StartTime: time.Date(0, 0, 0, 9, 0, 0, 0, time.UTC),
				EndTime:   time.Date(0, 0, 0, 10, 0, 0, 0, time.UTC),
			},
			//Thursday
			{
				Frequency:  ical.WeeklyFrequency,
				Day:        time.Thursday,
				StartTime:  time.Date(0, 0, 0, 9, 0, 0, 0, time.UTC),
				EndTime:    time.Date(0, 0, 0, 10, 0, 0, 0, time.UTC),
				Exceptions: []time.Time{time.Date(2025, 11, 27, 0, 0, 0, 0, time.UTC)},
			},
			//Every other Friday
			{
				Frequency: ical.BiWeeklyFrequency,
				Day:       time.Friday,
				StartTime: time.Date(0, 0, 0, 9, 0, 0, 0, time.UTC),
				EndTime:   time.Date(0, 0, 0, 10, 0, 0, 0, time.UTC),
			},
		},
		TimeZone: ical.TimeZone(timezones.US_Mountain),
		Location: "Conference Room A",
	})
	if err != nil {
		t.Fatalf("Failed to add multiple times a week event: %v", err)
	}

	err = cal.Save(fmt.Sprintf("./tmp/multiple_occurrences_calendar_output%s.ics", time.Now().Format("20060102150405")))
	if err != nil {
		t.Fatalf("Failed to save calendar: %v", err)
	}

	data, err := cal.Generate()
	if err != nil {
		t.Fatalf("Failed to generate calendar: %v", err)
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

func TestCalendarWithImage(t *testing.T) {
	cal := ical.Create("Image Calendar", "Calendar with event image.")
	image := ical.ImageFromURL("https://pkg.go.dev/static/shared/logo/go-blue.svg")
	err := cal.AddEvent(ical.Event{
		Title:       "Event with Image",
		Description: "This event has an associated image.",
		StartDate:   time.Date(2024, 9, 15, 10, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, 9, 15, 11, 0, 0, 0, time.UTC),
		TimeZone:    ical.TimeZone(timezones.US_Pacific),
		Image:       image,
	})
	if err != nil {
		t.Fatalf("Failed to add event with image: %v", err)
	}

	data, err := cal.Generate()
	if err != nil {
		t.Fatalf("Failed to generate calendar: %v", err)
	}

	if len(data) == 0 {
		t.Errorf("Generated calendar data is empty")
	}

	expectedImageString := "IMAGE;VALUE=URI;DISPLAY=BADGE;FMTTYPE=image/svg+xml:https://pkg.go.dev/static/shared/logo/go-blue.svg"
	if !strings.Contains(string(data), expectedImageString) {
		print(string(data))
		t.Errorf("Generated calendar data does not contain expected image string")
	}

	err = cal.Save(fmt.Sprintf("./tmp/image_calendar_output%s.ics", time.Now().Format("20060102150405")))
	if err != nil {
		t.Fatalf("Failed to save calendar: %v", err)
	}
}
