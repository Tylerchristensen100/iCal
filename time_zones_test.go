package ical

import (
	"testing"

	"github.com/Tylerchristensen100/iCal/timezones"
)

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

func TestTimeZoneValid(t *testing.T) {
	var tests = []struct {
		timeZone TimeZone
		expected bool
	}{
		{CentralTimeZone, true},
		{PacificTimeZone, true},
		{AlaskaTimeZone, true},
		{HawaiiTimeZone, true},
		{ArizonaTimeZone, true},
		{UTC, true},
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
	var tests = []struct {
		timeZone TimeZone
		expected []byte
	}{
		{CentralTimeZone, timezones.USCentral},
		{EasternTimeZone, timezones.USEastern},
		{MountainTimeZone, timezones.USMountain},
		{PacificTimeZone, timezones.USPacific},
		{AlaskaTimeZone, timezones.USAlaska},
		{ArizonaTimeZone, timezones.USArizona},
		{HawaiiTimeZone, timezones.USHawaii},
		{UTC, timezones.UTC},
		{"Invalid/Zone", timezones.UTC},
		{"", timezones.UTC},
	}

	for _, test := range tests {
		result := test.timeZone.ics()
		if string(result) != string(test.expected) {
			t.Errorf("TimeZoneToICS(%q) = %q; want %q", test.timeZone, result, test.expected)
		}
	}
}
