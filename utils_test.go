package ical

import (
	"testing"
	"time"
)

func TestDayOfWeekFromString(t *testing.T) {
	var tests = []struct {
		input    string
		expected time.Weekday
	}{
		{"monday", time.Monday},
		{"tuesday", time.Tuesday},
		{"wednesday", time.Wednesday},
		{"thursday", time.Thursday},
		{"friday", time.Friday},
		{"saturday", time.Saturday},
		{"sunday", time.Sunday},
	}
	for _, tt := range tests {
		result, err := DayOfWeekFromString(tt.input)
		if err != nil {
			t.Errorf("DayOfWeekFromString(%q) returned error: %v", tt.input, err)
		}
		if result != tt.expected {
			t.Errorf("DayOfWeekFromString(%q) returned %v, expected %v", tt.input, result, tt.expected)
		}
	}

	result, err := DayOfWeekFromString("invalid")
	if err == nil {
		t.Errorf("DayOfWeekFromString(\"invalid\") expected to return error but got none")
	}
	if result != 0 {
		t.Errorf("DayOfWeekFromString(\"invalid\") returned %v, expected 0", result)
	}
}
