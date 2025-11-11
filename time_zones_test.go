package ical

import "testing"

func TestTimeZoneToiCal(t *testing.T) {
	var tests = []struct {
		timeZone TimeZone
		expected string
	}{
		{CentralTimeZone, "America/Chicago"},
		{PacificTimeZone, "America/Los_Angeles"},
		{AlaskaTimeZone, "America/Anchorage"},
		{HawaiiTimeZone, "America/Honolulu"},
		{ArizonaTimeZone, "America/Phoenix"},
		{"", "UTC"},
		{"Invalid/Zone", "UTC"},
	}
	for _, test := range tests {
		result := test.timeZone.iCal()
		if result != test.expected {
			t.Errorf("TimeZoneToiCal(%q) = %q; want %q", test.timeZone, result, test.expected)
		}
	}
}
