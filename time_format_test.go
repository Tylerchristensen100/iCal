package ical

import (
	"testing"
	"time"
)

func TestStripDay(t *testing.T) {
	var tests = []struct {
		input    time.Time
		expected time.Time
	}{
		{
			input:    time.Date(2025, 7, 15, 14, 30, 0, 0, time.UTC),
			expected: time.Date(0, 1, 1, 14, 30, 0, 0, time.UTC),
		},
		{
			input:    time.Date(2025, 12, 25, 9, 15, 45, 0, time.UTC),
			expected: time.Date(0, 1, 1, 9, 15, 45, 0, time.UTC),
		},
	}

	for _, test := range tests {
		result := stripDay(test.input)
		if !result.Equal(test.expected) {
			t.Errorf("stripDay(%v) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestStripTime(t *testing.T) {
	var tests = []struct {
		input    time.Time
		expected time.Time
	}{
		{
			input:    time.Date(2025, 7, 15, 14, 30, 0, 0, time.UTC),
			expected: time.Date(2025, 7, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			input:    time.Date(2025, 12, 25, 9, 15, 45, 0, time.UTC),
			expected: time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, test := range tests {
		result := stripTime(test.input)
		if !result.Equal(test.expected) {
			t.Errorf("stripTime(%v) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestTimeToIcal(t *testing.T) {
	var tests = []struct {
		input    time.Time
		expected string
	}{
		{
			input:    time.Date(2025, 7, 15, 14, 30, 0, 0, time.UTC),
			expected: "20250715T143000",
		},
		{
			input:    time.Date(2025, 12, 25, 9, 15, 45, 0, time.UTC),
			expected: "20251225T091545",
		},
	}

	for _, test := range tests {
		result := timeToICal(test.input)
		if result != test.expected {
			t.Errorf("timeToICal(%v) = %v; want %v", test.input, result, test.expected)
		}
	}
}
