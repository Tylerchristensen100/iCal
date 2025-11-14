package ical

import (
	"strings"
	"testing"
	"time"
)

func TestJournal(t *testing.T) {
	var tests = []struct {
		name    string
		journal *Journal
		want    bool
	}{
		{
			name:    "Valid Journal",
			journal: mockJournal(),
			want:    true,
		},
		{
			name: "Invalid Journal - Empty Summary",
			journal: &Journal{
				Summary: "",
			},
			want: false,
		},
		{
			name: "Invalid Journal - Invalid Status",
			journal: &Journal{
				Summary: "Invalid Status Journal",
				Status:  "INVALID",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.journal.valid(); got != tt.want {
				t.Errorf("Journal.valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateJournal(t *testing.T) {
	journal := mockJournal()

	var builder strings.Builder
	err := journal.generate(&builder)
	if err != nil {
		t.Errorf("Error generating journal: %v", err)
	}

	if builder.Len() == 0 {
		t.Errorf("Generated journal is empty")
	}

	if !strings.Contains(builder.String(), "BEGIN:VJOURNAL") {
		t.Errorf("Generated journal missing BEGIN:VJOURNAL")
	}

	if !strings.Contains(builder.String(), "END:VJOURNAL") {
		t.Errorf("Generated journal missing END:VJOURNAL")
	}

	if !strings.Contains(builder.String(), "SUMMARY:Test Journal Entry") {
		t.Errorf("Generated journal missing SUMMARY")
	}

	if !strings.Contains(builder.String(), "ORGANIZER;CN=Test:mailto:test@example.com") {
		t.Errorf("Generated journal missing ORGANIZER")
	}
}

func TestValidJournal(t *testing.T) {
	journal := mockJournal()
	if !journal.valid() {
		t.Errorf("Expected valid journal, got invalid")
	}

	invalidJournal := &Journal{
		Summary: "",
	}
	if invalidJournal.valid() {
		t.Errorf("Expected invalid journal, got valid")
	}
}

func mockJournal() *Journal {
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	return &Journal{
		Summary:     "Test Journal Entry",
		Description: "This is a test journal entry for unit testing.",
		Organizer:   Participant{Name: "Test", Email: "test@example.com"},
		Status:      DraftJournal,
		StartDate:   &start,
	}

}
