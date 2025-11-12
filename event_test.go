package ical

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateEvent(t *testing.T) {
	cal := mockCalendar()
	event := mockEvent()
	err := cal.AddEvent(event)
	if err != nil {
		t.Fatalf("AddEvent() returned error: %v", err)
	}
	icalEvent, err := event.Generate()
	if err != nil {
		t.Errorf("Generate() returned error: %v", err)
	}
	if len(icalEvent) == 0 {
		t.Errorf("Generate() returned empty iCal event string")
	}

	if !strings.Contains(icalEvent, "BEGIN:VEVENT") || !strings.Contains(icalEvent, "END:VEVENT") {
		t.Errorf("Generated iCal event string is missing VEVENT boundaries")
	}

	if !strings.Contains(icalEvent, "ATTENDEE;") {
		t.Errorf("Generated iCal event string is missing ATTENDEE field")
	}
}

func TestHasRecurrences(t *testing.T) {
	cal := mockCalendar()
	event := mockEvent()
	err := cal.AddEvent(event)
	if err != nil {
		t.Fatalf("AddEvent() returned error: %v", err)
	}

	if !event.HasRecurrences() {
		t.Errorf("Expected event to have recurrences, but HasRecurrences() returned false")
	}

	nonRecurringEvent := Event{
		Title:     "Non-Recurring Event",
		StartDate: time.Now(),
		EndDate:   time.Now().Add(1 * time.Hour),
		TimeZone:  EasternTimeZone,
	}

	if nonRecurringEvent.HasRecurrences() {
		t.Errorf("Expected event to not have recurrences, but HasRecurrences() returned true")
	}
}

func TestEventUID(t *testing.T) {
	event := mockEvent()
	expectedUID := "Test_Event-Monday-Monday"
	if event.uid() != expectedUID {
		t.Errorf("Expected UID %s, got %s", expectedUID, event.uid())
	}
}

func TestConflictsWithEvent(t *testing.T) {
	cal := mockCalendar()
	event := mockEvent()
	err := cal.AddEvent(event)
	if err != nil {
		t.Fatalf("AddEvent() returned error: %v", err)
	}

	conflicts := cal.ListConflicts()
	if len(conflicts) == 0 {
		t.Errorf("Expected conflicts, but found %d", len(conflicts))
	}

	event1, event2 := conflicts[0], conflicts[1]
	conflict, tme := event1.ConflictsWith(event2)
	if !conflict {
		t.Errorf("Expected events to conflict, but ConflictsWith() returned false")
	}
	if tme.Weekday() != time.Monday {
		t.Errorf("Expected conflict day to be Monday, got %s", tme.Weekday())
	}

}

func TestCancelOnDate(t *testing.T) {
	event := mockEvent()
	exceptionDate := time.Date(2025, time.November, 24, 0, 0, 0, 0, time.UTC)
	err := event.CancelOnDate(exceptionDate)
	if err != nil {
		t.Fatalf("CancelOnDate() returned error: %v", err)
	}

	for i := range event.Recurrences {
		for _, ex := range event.Recurrences[i].Exceptions {
			if ex.Equal(exceptionDate) {
				t.Logf("Exception date %v successfully added to event exceptions", exceptionDate)
				return
			}
		}
	}
	t.Errorf("Expected exception date %v to be in event exceptions, but it was not found", exceptionDate)
}

func TestAddAttendee(t *testing.T) {
	event := mockEvent()
	var tests = []struct {
		email       string
		expectError bool
	}{
		{"test@example.com", false},
		{"invalid-email", true},
		{"", true},
	}

	for _, tt := range tests {
		err := event.AddAttendee(tt.email)
		if tt.expectError && err == nil {
			t.Errorf("Expected error for email %s, but got none", tt.email)
		}
		if !tt.expectError && err != nil {
			t.Errorf("Did not expect error for email %s, but got: %v", tt.email, err)
		}
		if !tt.expectError {
			found := false
			for _, attendee := range event.Attendees {
				if attendee == tt.email {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected attendee %s to be in event attendees, but it was not found", tt.email)
			}
		}
	}

}
func mockEvent() Event {
	startDate := time.Date(2025, time.November, 17, 9, 0, 0, 0, time.UTC)
	return Event{
		Title:     "Test Event",
		StartDate: startDate,
		EndDate:   startDate.Add(5 * time.Hour),
		TimeZone:  EasternTimeZone,
		Recurrences: []Recurrences{
			{
				Frequency: WeeklyFrequency,
				Day:       time.Monday,
				StartTime: time.Date(0, 0, 0, 9, 0, 0, 0, time.UTC),
				EndTime:   time.Date(0, 0, 0, 10, 0, 0, 0, time.UTC),
			},
		},
		Attendees: []string{"test@example.com"},
	}
}
