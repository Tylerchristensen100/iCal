package ical

import (
	"errors"
	"strings"
	"time"
)

// iCalendar VJOURNAL component
// https://icalendar.org/iCalendar-RFC-5545/3-6-3-journal-component.html

//  BEGIN:VJOURNAL
//  UID:19970901T130000Z-123405@example.com
//  DTSTAMP:19970901T130000Z
//  DTSTART;VALUE=DATE:19970317
//  SUMMARY:Staff meeting minutes
//  DESCRIPTION:1. Staff meeting: Participants include Joe\,
//    Lisa\, and Bob. Aurora project plans were reviewed.
//    There is currently no budget reserves for this project.
//    Lisa will escalate to management. Next meeting on Tuesday.\n
//   2. Telephone Conference: ABC Corp. sales representative
//    called to discuss new printer. Promised to get us a demo by
//    Friday.\n3. Henry Miller (Handsoff Insurance): Car was
//    totaled by tree. Is looking into a loaner car. 555-2323
//    (tel).
//  END:VJOURNAL

// VJournal component structure
type Journal struct {
	// OPTIONAL: The calendar date the journal entry is associated with.
	// This often uses the DATE value type (YYYYMMDD) if it's a daily entry,
	// but a DATE-TIME is also allowed.
	StartDate *time.Time // DTSTART

	// OPTIONAL: A brief summary or subject for the entry.
	Summary string

	// OPTIONAL: The full text of the journal entry. This is the main content.
	Description string

	// OPTIONAL: Current status of the journal entry.
	Status JournalStatus

	// OPTIONAL: Organizer's email/CN.
	Organizer Participant
}

type JournalStatus string

const (
	DraftJournal     JournalStatus = "DRAFT"
	FinalJournal     JournalStatus = "FINAL"
	CancelledJournal JournalStatus = "CANCELLED"
)

// TODO: Add to either Event or Calendar method (whatever it is used as)
func (j *Journal) generate(builder *strings.Builder) error {
	// UID
	// DTSTAMP
	return errors.New("Not implemented yet")
}

func (j *Journal) valid() bool {
	return true
}
