package ical

import (
	"strings"
	"testing"
	"time"

	"github.com/Tylerchristensen100/iCal/timezones"
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
		TimeZone:  TimeZone(timezones.US_Eastern),
	}

	if nonRecurringEvent.HasRecurrences() {
		t.Errorf("Expected event to not have recurrences, but HasRecurrences() returned true")
	}
}

func TestAddReminder(t *testing.T) {
	event := mockEvent()
	var tests = []struct {
		reminder    Reminder
		expectError bool
	}{
		{*mockReminder(), false},
		{Reminder{
			Action:      EmailReminderAction,
			Description: "Test Reminder",
			Trigger:     time.Duration(-15 * time.Minute),
			Attendees:   []Participant{},
		}, true},
		{Reminder{
			Action:      "INVALID",
			Description: "Test Reminder",
			Trigger:     time.Duration(-15 * time.Minute),
		}, true}, // Invalid: negative minutes
	}

	for _, tt := range tests {
		err := event.AddReminder(tt.reminder)
		if tt.expectError && err == nil {
			t.Errorf("Expected error for reminder %v, but got none", tt.reminder)
		}
		if !tt.expectError && err != nil {
			t.Errorf("Did not expect error for reminder %v, but got: %v", tt.reminder, err)
		}
		if !tt.expectError {
			found := false
			for _, reminder := range event.Reminders {
				if reminder.Description == tt.reminder.Description {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected reminder with description %s to be in event reminders, but it was not found", tt.reminder.Description)
			}
		}
	}
}

func TestEventUID(t *testing.T) {
	event := mockEvent()
	expectedUID := "Test_Event-Monday-Monday@iCal.go"
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
		err := event.AddAttendee(tt.email, tt.email)
		if tt.expectError && err == nil {
			t.Errorf("Expected error for email %s, but got none", tt.email)
		}
		if !tt.expectError && err != nil {
			t.Errorf("Did not expect error for email %s, but got: %v", tt.email, err)
		}
		if !tt.expectError {
			found := false
			for _, attendee := range event.Attendees {
				if attendee.Email == tt.email {
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

func TestAddOrganizer(t *testing.T) {
	event := mockEvent()
	err := event.AddOrganizer("Test Organizer", "test@example.com")
	if err != nil {
		t.Errorf("AddOrganizer() returned error: %v", err)
	}
	if event.Organizer.Name != "Test Organizer" || event.Organizer.Email != "test@example.com" {
		t.Errorf("Organizer not set correctly. Got: %+v", event.Organizer)
	}
}

func TestCleanDescription(t *testing.T) {
	rawDescription := "This is a test description with special characters: \n , ; \\ and more."
	cleanedDescription := cleanDescription(rawDescription)
	expectedDescription := "This is a test description with special characters:   , ,  a..."
	if cleanedDescription != expectedDescription {
		t.Errorf("Expected cleaned description to be:\n%s\nGot:\n%s", expectedDescription, cleanedDescription)
	}

	clean := "Short description."
	cleaned := cleanDescription(clean)
	if cleaned != clean {
		t.Errorf("Expected cleaned description to be unchanged:\n%s\nGot:\n%s", clean, cleaned)
	}

}

func TestConflictsWithSingleEvent(t *testing.T) {
	event1 := Event{
		Title:     "Event 1",
		StartDate: time.Date(2025, time.November, 17, 9, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, time.November, 17, 10, 0, 0, 0, time.UTC),
		TimeZone:  TimeZone(timezones.US_Eastern),
	}

	event2 := Event{
		Title:     "Event 2",
		StartDate: time.Date(2025, time.November, 17, 9, 30, 0, 0, time.UTC),
		EndDate:   time.Date(2025, time.November, 17, 10, 30, 0, 0, time.UTC),
		TimeZone:  TimeZone(timezones.US_Eastern),
	}

	conflict, _ := event1.ConflictsWith(&event2)
	if !conflict {
		t.Errorf("Expected events to conflict, but ConflictsWith() returned false")
	}
}

func TestConflictWithRecurringEvent(t *testing.T) {
	event1 := mockEvent()

	event2 := Event{
		Title:     "Event 2",
		StartDate: time.Date(2025, time.November, 24, 9, 30, 0, 0, time.UTC),
		EndDate:   time.Date(2025, time.November, 24, 10, 30, 0, 0, time.UTC),
		TimeZone:  TimeZone(timezones.US_Eastern),
	}
	conflict, conflictDate := event1.ConflictsWith(&event2)
	if !conflict {
		t.Errorf("Expected events to conflict, but ConflictsWith() returned false")
	}
	expectedDate := time.Date(2025, time.November, 24, 9, 0, 0, 0, time.UTC)
	if !conflictDate.Equal(expectedDate) {
		t.Errorf("Expected conflict date to be %v, got %v", expectedDate, conflictDate)
	}
}

func mockEvent() Event {
	startDate := time.Date(2025, time.November, 17, 9, 0, 0, 0, time.UTC)
	return Event{
		Title:     "Test Event",
		StartDate: startDate,
		EndDate:   startDate.Add(5 * time.Hour),
		TimeZone:  TimeZone(timezones.US_Eastern),
		Recurrences: []Recurrences{
			{
				Frequency: WeeklyFrequency,
				Day:       time.Monday,
				StartTime: time.Date(0, 0, 0, 9, 0, 0, 0, time.UTC),
				EndTime:   time.Date(0, 0, 0, 10, 0, 0, 0, time.UTC),
			},
		},
		Attendees: []Participant{{Name: "Test User", Email: "test@example.com"}},
	}
}
