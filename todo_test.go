package ical

import (
	"strings"
	"testing"
)

func TestTodo(t *testing.T) {
	var tests = []struct {
		name string
		todo *Todo
		want bool
	}{
		{
			name: "Valid Todo",
			todo: mockTodo(),
			want: true,
		},
		{
			name: "Invalid Todo - Empty Summary",
			todo: &Todo{
				Summary: "",
			},
			want: false,
		},
		{
			name: "Funky Priority",
			todo: &Todo{
				Summary:  "Funky Priority Todo",
				Priority: new(int),
			},
			want: false,
		},
		{
			name: "Invalid Recurrence",
			todo: &Todo{
				Summary: "Invalid Recurrence Todo",
				Recurrence: &Recurrences{
					Frequency: WeeklyFrequency,
					Day:       -1,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.todo.valid(); got != tt.want {
				t.Errorf("Todo.valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateTodo(t *testing.T) {
	todo := mockTodo()

	var builder strings.Builder
	err := todo.generate(&builder)
	if err != nil {
		t.Errorf("Error generating todo: %v", err)
	}

	if builder.Len() == 0 {
		t.Errorf("Generated todo is empty")
	}

	if !strings.Contains(builder.String(), "BEGIN:VTODO") {
		t.Errorf("Generated todo missing BEGIN:VTODO")
	}

	if !strings.Contains(builder.String(), "END:VTODO") {
		t.Errorf("Generated todo missing END:VTODO")
	}

	if !strings.Contains(builder.String(), "SUMMARY:Test To-Do Item") {
		t.Errorf("Generated todo missing SUMMARY")
	}
}

func TestGenerateTodoWithReminder(t *testing.T) {
	todo := mockTodo()
	reminder := mockReminder()
	todo.Alarms = []Reminder{*reminder}

	var builder strings.Builder
	err := todo.generate(&builder)
	if err != nil {
		t.Errorf("Error generating todo with reminder: %v", err)
	}

	if !strings.Contains(builder.String(), "BEGIN:VALARM") {
		t.Errorf("Generated todo missing BEGIN:VALARM for reminder")
	}

	if !strings.Contains(builder.String(), "END:VALARM") {
		t.Errorf("Generated todo missing END:VALARM for reminder")
	}

	if !strings.Contains(builder.String(), "ACTION:DISPLAY") {
		t.Errorf("Generated todo missing ACTION for reminder")
	}

	if !strings.Contains(builder.String(), "DESCRIPTION:Test") {
		t.Errorf("Generated todo missing DESCRIPTION for reminder")
	}
}

func TestValidTodo(t *testing.T) {
	todo := mockTodo()
	if !todo.valid() {
		t.Errorf("Expected valid todo, got invalid")
	}

	invalidTodo := &Todo{
		Summary: "",
	}
	if invalidTodo.valid() {
		t.Errorf("Expected invalid todo, got valid")
	}
}

func TestTodoStatus(t *testing.T) {
	var tests = []struct {
		status   TodoStatus
		expected bool
	}{
		{NeedsActionStatus, true},
		{CompletedStatus, true},
		{CancelledStatus, true},
		{"IN PROGRESS", false},
		{"INVALID_STATUS", false},
	}

	for _, tt := range tests {
		if tt.status.valid() != tt.expected {
			t.Errorf("TodoStatus.valid() = %v, want %v for status %v", tt.status.valid(), tt.expected, tt.status)
		}
	}
}

func mockTodo() *Todo {
	priority := 5
	completePercent := 3
	return &Todo{
		Summary:         "Test To-Do Item",
		Description:     "This is a test to-do item for unit testing.",
		Status:          NeedsActionStatus,
		Priority:        &priority,
		PercentComplete: &completePercent,
		Organizer:       Participant{Name: "Test", Email: "test@example.com"},
	}

}
