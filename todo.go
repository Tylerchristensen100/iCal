package ical

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// iCalendar VTODO component
// https://icalendar.org/iCalendar-RFC-5545/3-6-2-to-do-component.html
type Todo struct {
	Summary     string
	Description string

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

	// OPTIONAL: Start date/time of the To-Do.
	StartDate *time.Time

	// OPTIONAL: The Reminder (VALARM) components can be attached here.
	Alarms []Reminder

	// OPTIONAL: Organizer's email/CN.
	Organizer Participant

	Recurrence *Recurrences
}

// TODO: Add to either Event or Calendar method (whatever it is used as)
func (t *Todo) generate(builder *strings.Builder) error {
	if !t.valid() {
		return errors.New("Invalid Todo component")
	}
	builder.WriteString("BEGIN:VTODO" + lineBreak)
	builder.WriteString("UID:" + t.uid() + lineBreak)
	builder.WriteString("DTSTAMP:" + timeToICal(time.Now().UTC()) + lineBreak)
	builder.WriteString("SUMMARY:" + t.Summary + lineBreak)
	if t.Status != "" {
		builder.WriteString("STATUS:" + string(t.Status) + lineBreak)
	}
	if t.Due != nil {
		builder.WriteString("DUE:" + timeToICal(*t.Due) + lineBreak)
	}
	if t.Completed != nil {
		builder.WriteString("COMPLETED:" + timeToICal(*t.Completed) + lineBreak)
	}
	if t.Priority != nil {
		builder.WriteString(fmt.Sprintf("PRIORITY:%d%s", *t.Priority, lineBreak))
	}
	if t.PercentComplete != nil {
		builder.WriteString(fmt.Sprintf("PERCENT-COMPLETE:%d%s", *t.PercentComplete, lineBreak))
	}
	if t.Description != "" {
		description := cleanDescription(t.Description)
		builder.WriteString("DESCRIPTION:" + description + lineBreak)
	}
	if t.StartDate != nil {
		builder.WriteString("DTSTART:" + timeToICal(*t.StartDate) + lineBreak)
	}
	if t.Organizer.Email != "" {
		err := t.Organizer.generateOrganizer(builder)
		if err != nil {
			return err
		}
	}

	if len(t.Alarms) > 0 {
		for _, alarm := range t.Alarms {
			err := alarm.generate(builder)
			if err != nil {
				return err
			}
		}
	}

	if t.Recurrence != nil {
		builder.WriteString(generateRRULE(t.Recurrence.Frequency, t.Recurrence.Day, t.Recurrence.EndTime))
	}

	builder.WriteString("END:VTODO" + lineBreak)
	return nil
}

func (t *Todo) valid() bool {
	if t.Summary == "" {
		return false
	}

	if t.Recurrence != nil {
		if t.Recurrence.Frequency != "" && t.Recurrence.Day < 0 {
			return false
		}
	}

	if t.Status != "" && !t.Status.valid() {
		return false
	}
	return true
}

func (t *Todo) uid() string {
	return fmt.Sprintf("%s-%d@iCal.go", strings.ReplaceAll(t.Summary, " ", "_"), time.Now().Unix())
}

type TodoStatus string

const (
	NeedsActionStatus TodoStatus = "NEEDS-ACTION"
	CompletedStatus   TodoStatus = "COMPLETED"
	InProcessStatus   TodoStatus = "IN-PROCESS"
	CancelledStatus   TodoStatus = "CANCELLED"
)

func (s *TodoStatus) valid() bool {
	switch *s {
	case NeedsActionStatus, CompletedStatus, InProcessStatus, CancelledStatus:
		return true
	default:
		return false
	}
}
