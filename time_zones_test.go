package ical

import (
	"testing"

	"github.com/Tylerchristensen100/iCal/timezones"
)

func TestTimeZoneValid(t *testing.T) {
	var tests = []struct {
		timeZone TimeZone
		expected bool
	}{
		{TimeZone(timezones.America_Chicago), true},
		{TimeZone(timezones.America_Los_Angeles), true},
		{TimeZone(timezones.US_Alaska), true},
		{TimeZone(timezones.US_Hawaii), true},
		{TimeZone(timezones.US_Arizona), true},
		{TimeZone(timezones.UTC), true},
		{"", false},
		{"Invalid/Zone", false},
	}
	for _, test := range tests {
		result := test.timeZone.valid()
		if result != test.expected {
			t.Errorf("TimeZoneValid(%q) = %v; want %v", test.timeZone, result, test.expected)
		}
	}
}

func TestTimeZoneToICS(t *testing.T) {
	timezone := TimeZone(timezones.America_New_York)
	icsString, found := timezone.iCal()
	if !found {
		t.Errorf("Expected to find timezone definition for %q", timezone)
	}
	if icsString == "" {
		t.Errorf("Expected non-empty ICS string for timezone %q", timezone)
	}

	invalidTimezone := TimeZone("Invalid/Timezone")
	icsString, found = invalidTimezone.iCal()
	if found {
		t.Errorf("Did not expect to find timezone definition for %q", invalidTimezone)
	}
	if icsString != string(timezones.UTC) {
		t.Errorf("Expected ICS string for invalid timezone %q to be UTC, got %q", invalidTimezone, icsString)
	}
}
