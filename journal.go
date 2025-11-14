package ical

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// iCalendar VJOURNAL component
// https://icalendar.org/iCalendar-RFC-5545/3-6-3-journal-component.html
type Journal struct {

	// REQUIRED: Short summary or title of the journal entry
	Summary string

	// REQUIRED: Detailed description of the journal entry
	Description string

	// OPTIONAL: Current status of the journal entry.
	//
	// Possible Values: DraftJournal, FinalJournal, CancelledJournal
	Status JournalStatus

	// OPTIONAL: Start date of the journal entry.
	StartDate *time.Time

	// OPTIONAL: Organizer's name and email
	Organizer Participant
}

type JournalStatus string

func (j *Journal) generate(builder *strings.Builder) error {
	if !j.valid() {
		return errors.New("invalid journal entry")
	}
	builder.WriteString("BEGIN:VJOURNAL" + lineBreak)

	builder.WriteString("UID:" + j.uid() + lineBreak)
	builder.WriteString("DTSTAMP:" + timeToICal(time.Now().UTC()) + lineBreak)

	if j.StartDate != nil {
		builder.WriteString("DTSTART;VALUE=DATE:" + j.StartDate.Format("20060102") + lineBreak)
	}

	builder.WriteString("SUMMARY:" + j.Summary + lineBreak)
	builder.WriteString("DESCRIPTION:" + cleanDescription(j.Description) + lineBreak)

	if j.Organizer.Email != "" {
		builder.WriteString("ORGANIZER;CN=" + j.Organizer.Name + ":mailto:" + j.Organizer.Email + lineBreak)
	}

	if j.Status != "" {
		builder.WriteString("STATUS:" + string(j.Status) + lineBreak)
	}

	builder.WriteString("END:VJOURNAL" + lineBreak)
	return nil
}

func (j *Journal) valid() bool {
	if j.Status != "" && !j.Status.valid() {
		return false
	}

	if j.Summary == "" {
		return false
	}

	if j.Description == "" {
		return false
	}

	return true
}

func (j *Journal) uid() string {
	return fmt.Sprintf("%s-%d@iCal.go", strings.ReplaceAll(j.Summary, " ", "_"), time.Now().Unix())
}

const (
	DraftJournal     JournalStatus = "DRAFT"
	FinalJournal     JournalStatus = "FINAL"
	CancelledJournal JournalStatus = "CANCELLED"
)

func (s *JournalStatus) valid() bool {
	switch *s {
	case DraftJournal, FinalJournal, CancelledJournal:
		return true
	default:
		return false
	}
}
