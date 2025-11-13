package ical

import (
	"strings"
	"testing"
	"time"

	"github.com/Tylerchristensen100/iCal/timezones"
)

func TestAddEvent(t *testing.T) {
	var tests = []struct {
		event   Event
		success bool
	}{
		{Event{Title: "Valid Event",
			StartDate: time.Now(),
			EndDate:   time.Now().Add(1 * time.Hour),
			TimeZone:  TimeZone(timezones.US_Eastern),
		}, true},
		{Event{Title: "", StartDate: time.Now(), EndDate: time.Now().Add(1 * time.Hour), TimeZone: TimeZone(timezones.US_Eastern)}, false},              // Invalid: empty title
		{Event{Title: "Invalid Event", StartDate: time.Now().Add(2 * time.Hour), EndDate: time.Now(), TimeZone: TimeZone(timezones.US_Eastern)}, false}, // Invalid: end before start
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

func TestGenerateTimeZones(t *testing.T) {
	builder := strings.Builder{}

	validtz := TimeZone(timezones.US_Central)

	cal := mockCalendar(validtz)
	cal.generateTimeZones(&builder)
	result := builder.String()
	expectedTimeZones, found := timezones.Get(timezones.US_Central)
	if !found || !strings.Contains(result, string(expectedTimeZones)) {
		t.Errorf("Expected time zone definition not found in generated iCal data")
	}
	builder.Reset()

	invalidtz := TimeZone("Invalid/Timezone")
	cal = mockCalendar(invalidtz)
	cal.generateTimeZones(&builder)
	result = builder.String()
	expectedTimeZones, found = timezones.Get(timezones.UTC)
	if !found && !strings.Contains(result, string(expectedTimeZones)) {
		t.Errorf("Expected UTC time zone definition not found in generated iCal data for invalid timezone")
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

func mockCalendar(tz ...TimeZone) *Calendar {
	var zone TimeZone = TimeZone(timezones.US_Eastern)
	if tz != nil {
		zone = tz[0]
	}
	cal := Create("Test Calendar", "A calendar for testing")
	cal.AddEvent(Event{Title: "Valid Event",
		StartDate: time.Now(),
		EndDate:   time.Now().Add(1 * time.Hour),
		TimeZone:  zone,
	})
	cal.AddEvent(mockEvent())
	return cal
}
