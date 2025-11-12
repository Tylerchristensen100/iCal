package ical

import (
	"testing"
	"time"
)

func TestGenerateRecurrences(t *testing.T) {
	var tests = []struct {
		recurrences Recurrences
		success     bool
	}{
		{mockRecurrence(), true},
		{Recurrences{Frequency: WeeklyFrequency}, false},
	}
	for _, tt := range tests {
		s, err := tt.recurrences.Generate(time.Now(), time.Now().Add(7*24*time.Hour), EasternTimeZone)
		if err != nil && tt.success {
			t.Errorf("Recurrences.Generate() returned error: %v", err)
		}
		if err == nil && !tt.success {
			t.Errorf("Recurrences.Generate() expected to fail but succeeded, output: %s", s)
		}
	}
}

func TestGenerateRecurrencesWithException(t *testing.T) {
	var tests = []struct {
		recurrences Recurrences
		success     bool
	}{
		{mockRecurrence(), true},
		{Recurrences{Frequency: WeeklyFrequency, Exceptions: []time.Time{time.Now().Add(24 * time.Hour)}}, false},
		{Recurrences{Frequency: DailyFrequency, Exceptions: []time.Time{time.Now().Add(24 * 4 * time.Hour)}}, false},
	}
	for _, tt := range tests {
		s, err := tt.recurrences.Generate(time.Now(), time.Now().Add(7*24*time.Hour), EasternTimeZone)
		if err != nil && tt.success {
			t.Errorf("Recurrences.Generate() returned error: %v", err)
		}
		if err == nil && !tt.success {
			t.Errorf("Recurrences.Generate() expected to fail but succeeded, output: %s", s)
		}
	}
}

func TestConflictsWithRecurrences(t *testing.T) {
	r := mockRecurrence()
	other := Recurrences{
		Frequency: WeeklyFrequency,
		Day:       time.Monday,
		StartTime: time.Date(0, 0, 0, 9, 30, 0, 0, time.UTC),
		EndTime:   time.Date(0, 0, 0, 10, 30, 0, 0, time.UTC),
	}
	conflict, tme := r.ConflictsWith(other)
	if !conflict {
		t.Errorf("Expected conflict between recurrences but got none")
	}
	if tme.Hour() != 9 || tme.Minute() != 00 {
		t.Errorf("Expected conflict at 9:00 but got %d:%d", tme.Hour(), tme.Minute())
	}

	nonConflicting := Recurrences{
		Frequency: WeeklyFrequency,
		Day:       time.Tuesday,
		StartTime: time.Date(0, 0, 0, 9, 30, 0, 0, time.UTC),
		EndTime:   time.Date(0, 0, 0, 10, 30, 0, 0, time.UTC),
	}
	conflict, _ = r.ConflictsWith(nonConflicting)
	if conflict {
		t.Errorf("Expected no conflict between recurrences but got one")
	}

}
func TestFindStartDate(t *testing.T) {
	var tests = []struct {
		startDate time.Time
		startTime time.Time
		day       time.Weekday
		expected  time.Time
	}{
		{
			startDate: time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC), // Monday
			startTime: time.Date(0, 0, 0, 9, 0, 0, 0, time.UTC),
			day:       time.Wednesday,
			expected:  time.Date(2024, 7, 3, 9, 0, 0, 0, time.UTC),
		},
		{
			startDate: time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC), // Monday
			startTime: time.Date(0, 0, 0, 9, 0, 0, 0, time.UTC),
			day:       time.Monday,
			expected:  time.Date(2024, 7, 1, 9, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		date, err := findStartDate(tt.startDate, tt.day, tt.startTime)
		if err != nil {
			t.Errorf("findStartDate() returned error: %v", err)
		}
		if !date.Equal(tt.expected) {
			t.Errorf("findStartDate() returned %v, expected %v", date, tt.expected)
		}
	}
}

func TestFindEndDate(t *testing.T) {
	var tests = []struct {
		startDate time.Time
		startTime time.Time
		day       time.Weekday
		expected  time.Time
	}{
		{
			startDate: time.Date(2024, 7, 10, 0, 0, 0, 0, time.UTC), // Wednesday
			startTime: time.Date(0, 0, 0, 17, 0, 0, 0, time.UTC),
			day:       time.Monday,
			expected:  time.Date(2024, 7, 8, 17, 0, 0, 0, time.UTC),
		},
		{
			startDate: time.Date(2024, 7, 8, 0, 0, 0, 0, time.UTC), // Monday
			startTime: time.Date(0, 0, 0, 17, 0, 0, 0, time.UTC),
			day:       time.Monday,
			expected:  time.Date(2024, 7, 8, 17, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		date, err := findEndDate(tt.startDate, tt.day, tt.startTime)
		if err != nil {
			t.Errorf("findEndDate() returned error: %v", err)
		}
		if !date.Equal(tt.expected) {
			t.Errorf("findEndDate() returned %v, expected %v", date, tt.expected)
		}
	}
}

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

func TestWeekdayToICal(t *testing.T) {
	var tests = []struct {
		input    time.Weekday
		expected string
	}{
		{time.Monday, "MO"},
		{time.Tuesday, "TU"},
		{time.Wednesday, "WE"},
		{time.Thursday, "TH"},
		{time.Friday, "FR"},
		{time.Saturday, "SA"},
		{time.Sunday, "SU"},
	}
	for _, tt := range tests {
		result := weekdayToICal(tt.input)
		if result != tt.expected {
			t.Errorf("weekdayToICal(%v) returned %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestRecurrenceUID(t *testing.T) {
	rec := mockRecurrence()
	expectedUID := "WEEKLY-09_00-10_00"
	if rec.uid() != expectedUID {
		t.Errorf("Expected UID %s, got %s", expectedUID, rec.uid())
	}
}

func mockRecurrence() Recurrences {
	return Recurrences{
		Frequency: WeeklyFrequency,
		Day:       time.Monday,
		StartTime: time.Date(0, 0, 0, 9, 0, 0, 0, time.UTC),
		EndTime:   time.Date(0, 0, 0, 10, 0, 0, 0, time.UTC),
	}
}
