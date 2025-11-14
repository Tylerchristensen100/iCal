package ical

import (
	"errors"
	"strings"
	"time"
)

// iCalendar VTODO component
// https://icalendar.org/iCalendar-RFC-5545/3-6-2-to-do-component.html

//  BEGIN:VTODO
//  UID:20070313T123432Z-456553@example.com
//  DTSTAMP:20070313T123432Z
//  DUE;VALUE=DATE:20070501
//  SUMMARY:Submit Quebec Income Tax Return for 2006
//  CLASS:CONFIDENTIAL
//  CATEGORIES:FAMILY,FINANCE
//  STATUS:NEEDS-ACTION
//  END:VTODO

type Todo struct {

	// REQUIRED: A brief summary or subject for the To-Do.
	Summary string

	// OPTIONAL: Date/Time the To-Do is due. Use a pointer for optionality.
	Due *time.Time

	// OPTIONAL: Date/Time the To-Do was actually completed.
	Completed *time.Time

	// OPTIONAL: Current status of the To-Do.
	Status TodoStatus

	// OPTIONAL: Priority (0-9).
	Priority *int

	// OPTIONAL: Percentage of the To-Do completed (0-100).
	PercentComplete *int

	// OPTIONAL: Detailed description.
	Description string

	// OPTIONAL: Start date/time of the To-Do.
	DTStart *time.Time

	// OPTIONAL: The Reminder (VALARM) components can be attached here.
	// Assuming you have a Reminder struct from the previous conversation.
	Alarms []Reminder

	// OPTIONAL: Organizer's email/CN.
	Organizer Participant

	// OPTIONAL: Recurrence Rule. You'd typically use a separate struct/string for this.
	RRule string
}

type TodoStatus string

const (
	NeedsActionStatus TodoStatus = "NEEDS-ACTION"
	CompletedStatus   TodoStatus = "COMPLETED"
	InProcessStatus   TodoStatus = "IN-PROCESS"
	CancelledStatus   TodoStatus = "CANCELLED"
)

// TODO: Add to either Event or Calendar method (whatever it is used as)
func (t *Todo) generate(builder *strings.Builder) error {
	// UID
	// DTSTAMP
	// SUMMARY
	// DUE
	// COMPLETED
	// STATUS
	// PRIORITY
	// PERCENT-COMPLETE
	// DESCRIPTION
	// DTSTART
	// ORGANIZER
	// RRULE
	// VALARM (for each alarm in Alarms)
	return errors.New("Not implemented yet")

}

func (t *Todo) valid() bool {
	return true
}
