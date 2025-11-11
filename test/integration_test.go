package test

import (
	"os"
	"strings"
	"testing"
	"time"

	ical "github.com/Tylerchristensen100/iCal"
)

func TestIntegration(t *testing.T) {
	cal := ical.Create("Test Calendar", "This is a test calendar.")
	cal.AddEvent(ical.Event{
		Title:       "Weekly Meeting",
		Description: "Recurring Weekly",
		StartDate:   time.Date(2024, 7, 1, 9, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2026, 7, 1, 10, 0, 0, 0, time.UTC),
		Recurrences: []ical.Recurrences{{
			Frequency: ical.WeeklyFrequency,
			Day:       time.Monday,
			StartTime: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),
			EndTime:   time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
		},
		},
	})

	cal.AddEvent(ical.Event{
		Title:       "Monthly Meeting",
		Description: "Recurring Monthly",
		StartDate:   time.Date(2024, 7, 3, 14, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2026, 7, 3, 15, 0, 0, 0, time.UTC),
		Recurrences: []ical.Recurrences{{
			Frequency:  ical.MonthlyFrequency,
			Day:        time.Wednesday,
			StartTime:  time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC),
			EndTime:    time.Date(0, 1, 1, 15, 0, 0, 0, time.UTC),
			Exceptions: []time.Time{time.Date(2025, 12, 3, 14, 0, 0, 0, time.UTC)},
		},
		},
	})

	calendar, err := cal.Generate()
	if err != nil {
		t.Fatalf("Failed to generate calendar: %v", err)
	}

	expectedSubstring := "BEGIN:VCALENDAR"
	if !strings.Contains(calendar, expectedSubstring) {
		t.Errorf("Generated calendar does not contain expected substring: %s", expectedSubstring)
	}

	expectedEvent := "SUMMARY:Weekly Meeting"
	if !strings.Contains(calendar, expectedEvent) {
		t.Errorf("Generated calendar does not contain expected event: %s", expectedEvent)
	}

	expectedRecurrence := "RRULE:FREQ=WEEKLY;BYDAY=MO;"
	if !strings.Contains(calendar, expectedRecurrence) {
		t.Errorf("Generated calendar does not contain expected recurrence rule: %s", expectedRecurrence)
	}

	err = os.WriteFile("result_integration_test.ics", []byte(calendar), 0644)
	if err != nil {
		t.Fatalf("Failed to write calendar to file: %v", err)
	}
}
