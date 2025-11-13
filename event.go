package ical

import (
	"fmt"
	"strings"
	"time"
)

type Event struct {
	Title       string
	Description string
	Location    string
	Organizer   *Participant
	StartDate   time.Time
	EndDate     time.Time
	TimeZone    TimeZone

	// Recurrence rules for the event
	// If empty, the event is non-recurring
	Recurrences []Recurrences

	// Everyone invited to the event
	// Should be in email format
	Attendees []Participant
}

// Generate creates the iCal formatted string for the event.
func (e *Event) Generate() (string, error) {
	var builder strings.Builder

	if len(e.Recurrences) > 0 {
		// Recurring
		for _, rec := range e.Recurrences {
			recRule, err := rec.Generate(e.StartDate, e.EndDate, e.TimeZone)
			if err != nil {
				return "", err
			}

			builder.WriteString("BEGIN:VEVENT" + lineBreak)
			builder.WriteString("UID:" + rec.uid() + lineBreak)
			builder.WriteString(recRule + lineBreak)
			err = e.buildEventDetails(&builder)
			if err != nil {
				return "", err
			}

			builder.WriteString("END:VEVENT" + lineBreak)
		}
	} else {
		// Non-recurring
		builder.WriteString("BEGIN:VEVENT" + lineBreak)
		builder.WriteString("UID:" + e.uid() + lineBreak)

		builder.WriteString(fmt.Sprintf("DTSTART;TZID=%s:%s", e.TimeZone.ID(), timeToICal(e.StartDate)) + lineBreak)
		builder.WriteString(fmt.Sprintf("DTEND;TZID=%s:%s", e.TimeZone.ID(), timeToICal(e.EndDate)) + lineBreak)
		err := e.buildEventDetails(&builder)
		if err != nil {
			return "", err
		}

		builder.WriteString("END:VEVENT" + lineBreak)
	}

	return builder.String(), nil
}

// Adds a cancellation for the event on the specified date.
//
// If the event does not recur on that day, returns an error.
//
// Strips out the time component of the cancelDate to only consider the date.
func (e *Event) CancelOnDate(cancelDate time.Time) error {
	dayOfWeek := cancelDate.Weekday()

	date := time.Date(cancelDate.Year(), cancelDate.Month(), cancelDate.Day(),
		0, 0, 0, 0, cancelDate.Location())

	for i := range e.Recurrences {
		if e.Recurrences[i].Day == dayOfWeek {
			e.Recurrences[i].Exceptions = append(e.Recurrences[i].Exceptions, date)
			return nil
		}
	}
	return ErrNoRecurrenceFound
}

func (e *Event) AddAttendee(name, email string) error {
	if !validateEmail(email) {
		return ErrInvalidEmail
	}

	e.Attendees = append(e.Attendees, Participant{Name: name, Email: email})
	return nil
}

func (e *Event) AddOrganizer(name, email string) error {
	if !validateEmail(email) {
		return ErrInvalidEmail
	}

	e.Organizer = &Participant{Name: name, Email: email}
	return nil
}

func (e *Event) uid() string {
	return fmt.Sprintf("%s-%s-%s", strings.ReplaceAll(e.Title, " ", "_"),
		e.StartDate.Weekday(), e.EndDate.Weekday())
}

func (e *Event) buildEventDetails(builder *strings.Builder) error {
	builder.WriteString(fmt.Sprintf("DTSTAMP:%s", fmt.Sprintf("%sZ", timeToICal(time.Now().UTC()))) + lineBreak)
	builder.WriteString("SUMMARY:" + e.Title + lineBreak)
	builder.WriteString("LOCATION:" + e.Location + lineBreak)

	if e.Organizer != nil {
		err := e.Organizer.generateOrganizer(builder)
		if err != nil {
			return err
		}
	}

	if e.Description != "" {
		description := cleanDescription(e.Description)
		builder.WriteString("DESCRIPTION:" + description + lineBreak)
	}
	for _, attendee := range e.Attendees {
		err := attendee.generate(builder)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Event) HasRecurrences() bool {
	return len(e.Recurrences) > 0
}
func (e *Event) ConflictsWith(other *Event) (bool, time.Time) {
	//Both are recurring events
	if e.HasRecurrences() && other.HasRecurrences() {
		for _, rec := range e.Recurrences {
			for _, otherRec := range other.Recurrences {
				if conflicts, tme := rec.ConflictsWith(otherRec); conflicts {
					// Find the date of the conflict
					occurrences := rec.Occurrences(e.StartDate, e.EndDate)
					otherOccurrences := otherRec.Occurrences(other.StartDate, other.EndDate)
					for _, occ := range occurrences {
						for _, otherOcc := range otherOccurrences {
							if occ.Year() == otherOcc.Year() && occ.Month() == otherOcc.Month() && occ.Day() == otherOcc.Day() {
								// The Day and Time of the conflict
								return true, time.Date(occ.Year(), occ.Month(), occ.Day(), tme.Hour(), tme.Minute(), tme.Second(), 0, occ.Location())
							}
						}
					}
				}
			}
		}
	}

	// Single event vs single event
	if !e.HasRecurrences() && !other.HasRecurrences() {
		if e.StartDate.Before(other.EndDate) && e.EndDate.After(other.StartDate) {
			return true, time.Date(e.StartDate.Year(), e.StartDate.Month(), e.StartDate.Day(),
				e.StartDate.Hour(), e.StartDate.Minute(), e.StartDate.Second(), 0, e.StartDate.Location())
		}
	}

	// One is single, one is recurring
	var recurring *Event
	var single *Event

	if e.HasRecurrences() && !other.HasRecurrences() {
		recurring = e
		single = other
	} else if !e.HasRecurrences() && other.HasRecurrences() {
		recurring = other
		single = e
	}

	if recurring != nil && single != nil {

		singleDay := single.StartDate.Weekday()

		for _, rec := range recurring.Recurrences {
			if rec.Day == singleDay {

				singleStart := single.StartDate
				singleEnd := single.EndDate

				recStart := time.Date(singleStart.Year(), singleStart.Month(), singleStart.Day(),
					rec.StartTime.Hour(), rec.StartTime.Minute(), rec.StartTime.Second(), 0, singleStart.Location())
				recEnd := time.Date(singleStart.Year(), singleStart.Month(), singleStart.Day(),
					rec.EndTime.Hour(), rec.EndTime.Minute(), rec.EndTime.Second(), 0, singleStart.Location())

				if singleStart.Before(recEnd) && singleEnd.After(recStart) {
					//Find the date of the conflict
					return true, time.Date(singleStart.Year(), singleStart.Month(), singleStart.Day(),
						rec.StartTime.Hour(), rec.StartTime.Minute(), rec.StartTime.Second(), 0, singleStart.Location())
				}
			}
		}
	}

	return false, time.Time{}
}

func (e *Event) Valid() bool {
	if e.Title == "" {
		return false
	}
	if e.EndDate.Before(e.StartDate) || e.EndDate.Equal(e.StartDate) {
		return false
	}
	return true
}

func cleanDescription(desc string) string {
	// Truncate description to 75 characters (including 'DESCRIPTION:' prefix)
	// https://icalendar.org/iCalendar-RFC-5545/3-1-content-lines.html
	var description string
	description = strings.ReplaceAll(desc, ";", ",")
	description = strings.ReplaceAll(description, "\\", "")
	if len(description) > 63 {
		description = description[:60] + "..."
	}
	return strings.ReplaceAll(description, "\n", " ")
}

func timeToICal(t time.Time) string {
	return t.Format(iCalTimeLayout)
}

func stripDay(t time.Time) time.Time {
	return time.Date(0, 1, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}
