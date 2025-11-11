package ical

import (
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
	conflict, day := event1.ConflictsWith(event2)
	if !conflict {
		t.Errorf("Expected events to conflict, but ConflictsWith() returned false")
	}
	if day != time.Monday {
		t.Errorf("Expected conflict day to be Monday, got %s", day)
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
	}
}
