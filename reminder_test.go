package ical

import (
	"strings"
	"testing"
	"time"
)

func TestString(t *testing.T) {
	reminder := mockReminder()
	result := reminder.String()
	if !strings.Contains(result, "Test Reminder") {
		t.Errorf("Expected description to be in string representation")
	}
}

func TestGenerateReminder(t *testing.T) {
	reminder := mockReminder()
	var builder strings.Builder
	err := reminder.generate(&builder)
	if err != nil {
		t.Errorf("Unexpected error generating reminder: %v", err)
	}
	result := builder.String()
	if !strings.Contains(result, "BEGIN:VALARM") {
		t.Errorf("Expected VALARM beginning in generated data")
	}
	if !strings.Contains(result, "ACTION:DISPLAY") {
		t.Errorf("Expected ACTION:DISPLAY in generated data")
	}
	if !strings.Contains(result, "DESCRIPTION:Test Reminder") {
		t.Errorf("Expected DESCRIPTION in generated data")
	}
	if !strings.Contains(result, "TRIGGER:-PT15M") {
		t.Errorf("Expected TRIGGER:-PT15M in generated data")
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{duration: -15 * time.Minute, expected: "-PT15M"},
		{duration: 30 * time.Minute, expected: "PT30M"},
		{duration: -2 * time.Hour, expected: "-PT2H"},
		{duration: 1*time.Hour + 30*time.Minute, expected: "PT1H30M"},
		{duration: -25 * time.Hour, expected: "-P1DT1H"},
		{duration: 3*24*time.Hour + 4*time.Hour + 5*time.Minute + 6*time.Second, expected: "P3DT4H5M6S"},
	}

	for _, test := range tests {
		result := formatDurationAsTrigger(test.duration)
		if result != test.expected {
			t.Errorf("For duration %v, expected %s but got %s", test.duration, test.expected, result)
		}
	}
}

func TestInvalidReminder(t *testing.T) {
	invalidReminders := []Reminder{
		{Action: "INVALID", Description: "Desc", Trigger: -15 * time.Minute},
		{Action: EmailReminderAction, Description: "Desc", Trigger: -15 * time.Minute, Attendees: []Participant{}},
		{Action: DisplayReminderAction, Description: "", Trigger: -15 * time.Minute},
	}

	for _, reminder := range invalidReminders {
		if reminder.valid() {
			t.Errorf("Expected reminder to be invalid: %+v", reminder)
		}
	}
}

func TestValidReminders(t *testing.T) {
	validReminders := []Reminder{
		{Action: DisplayReminderAction, Description: "Desc", Trigger: -15 * time.Minute},
		{Action: EmailReminderAction, Description: "Desc", Trigger: -15 * time.Minute, Attendees: []Participant{{Name: "John", Email: "john@example.com"}}},
	}

	for _, reminder := range validReminders {
		if !reminder.valid() {
			t.Errorf("Expected reminder to be valid: %+v", reminder)
		}
	}
}

func TestAddDuration(t *testing.T) {
	reminder := mockReminder()
	reminder.addTrigger(-30 * time.Minute)
	if reminder.Trigger != -30*time.Minute {
		t.Errorf("Expected trigger to be updated to -30 minutes, got %v", reminder.Trigger)
	}
}

func mockReminder() *Reminder {
	return &Reminder{
		Action:      DisplayReminderAction,
		Description: "Test Reminder",
		Trigger:     time.Duration(-15 * time.Minute),
	}
}
