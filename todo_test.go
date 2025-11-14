package ical

import "testing"

func TestTodo(t *testing.T) {
	t.Errorf("Not implemented yet")
}

func TestGenerateTodo(t *testing.T) {
	t.Errorf("Not implemented yet")
}

func TestValidTodo(t *testing.T) {
	t.Errorf("Not implemented yet")
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
