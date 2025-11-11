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

func TestConflictsWithRecurrences(t *testing.T) {
	r := mockRecurrence()
	other := Recurrences{
		Frequency: WeeklyFrequency,
		Day:       time.Monday,
		StartTime: time.Date(0, 0, 0, 9, 30, 0, 0, time.UTC),
		EndTime:   time.Date(0, 0, 0, 10, 30, 0, 0, time.UTC),
	}
	conflict, day := r.ConflictsWith(other)
	if !conflict {
		t.Errorf("Expected conflict between recurrences but got none")
	}
	if day != time.Monday {
		t.Errorf("Expected conflict on Monday but got %s", day)
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

func mockRecurrence() Recurrences {
	return Recurrences{
		Frequency: WeeklyFrequency,
		Day:       time.Monday,
		StartTime: time.Date(0, 0, 0, 9, 0, 0, 0, time.UTC),
		EndTime:   time.Date(0, 0, 0, 10, 0, 0, 0, time.UTC),
	}
}
