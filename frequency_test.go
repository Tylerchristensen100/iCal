package ical

import "testing"

func TestFrequencyValid(t *testing.T) {
	var tests = []struct {
		freq     Frequency
		expected bool
	}{
		{"DAILY", true},
		{"WEEKLY", true},
		{"MONTHLY", true},
		{"YEARLY", true},
		{"HOURLY", false},
		{"", false},
		{"INVALID", false},
	}

	for _, test := range tests {
		f := Frequency(test.freq)
		result := f.Valid()
		if result != test.expected {
			t.Errorf("Frequency.Valid() for '%s' = %v; want %v", test.freq, result, test.expected)
		}
	}
}
