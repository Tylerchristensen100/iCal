package ical

import (
	"testing"
	"time"
)

func TestAddEvent(t *testing.T) {
	var tests = []struct {
		event   Event
		success bool
	}{
		{Event{Title: "Valid Event",
			StartDate: time.Now(),
			EndDate:   time.Now().Add(1 * time.Hour),
			TimeZone:  EasternTimeZone,
		}, true},
		{Event{Title: "", StartDate: time.Now(), EndDate: time.Now().Add(1 * time.Hour), TimeZone: EasternTimeZone}, false},              // Invalid: empty title
		{Event{Title: "Invalid Event", StartDate: time.Now().Add(2 * time.Hour), EndDate: time.Now(), TimeZone: EasternTimeZone}, false}, // Invalid: end before start
	}
	cal := mockCalendar()
	for _, tt := range tests {
		err := cal.AddEvent(tt.event)
		if (err == nil) != tt.success {
			t.Errorf("AddEvent(%v) = %v, want success %v", tt.event, err, tt.success)
		}
	}
}

func TestGenerateCalendar(t *testing.T) {
	cal := mockCalendar()
	ical, err := cal.Generate()
	if err != nil {
		t.Errorf("Generate() returned error: %v", err)
	}
	if len(ical) == 0 {
		t.Errorf("Generate() returned empty iCal string")
	}
}

func TestListConflicts(t *testing.T) {
	cal := mockCalendar()
	conflicts := cal.ListConflicts()
	if len(conflicts) != 0 {
		t.Errorf("Expected no conflicts, but found %d", len(conflicts))
	}

}

func TestResolveConflicts(t *testing.T) {
	cal := mockCalendar()

	cal.ResolveConflicts(func(event1, event2 *Event, day time.Time) {
		t.Errorf("Unexpected conflict between '%s' and '%s' on %s", event1.Title, event2.Title, day)
	})

}

func mockCalendar() *Calendar {
	cal := Create("Test Calendar", "A calendar for testing")
	cal.AddEvent(Event{Title: "Valid Event",
		StartDate: time.Now(),
		EndDate:   time.Now().Add(1 * time.Hour),
		TimeZone:  EasternTimeZone,
	})
	cal.AddEvent(mockEvent())
	return cal
}
