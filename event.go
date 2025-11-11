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
	StartDate   time.Time
	EndDate     time.Time
	TimeZone    TimeZone
	Recurrences []Recurrences
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
			e.buildEventDetails(&builder)
			builder.WriteString("END:VEVENT" + lineBreak)
		}
	} else {
		// Non-recurring
		builder.WriteString("BEGIN:VEVENT" + lineBreak)
		builder.WriteString("UID:" + e.uid() + lineBreak)

		builder.WriteString(fmt.Sprintf("DTSTART;TZID=%s:%s", e.TimeZone.iCal(), timeToICal(e.StartDate)) + lineBreak)
		builder.WriteString(fmt.Sprintf("DTEND;TZID=%s:%s", e.TimeZone.iCal(), timeToICal(e.EndDate)) + lineBreak)
		e.buildEventDetails(&builder)
		builder.WriteString("END:VEVENT" + lineBreak)
	}

	return builder.String(), nil
}

func (e *Event) uid() string {
	return fmt.Sprintf("%s-%s-%s", strings.ReplaceAll(e.Title, " ", "_"),
		e.StartDate.Weekday(), e.EndDate.Weekday())
}

func (e *Event) buildEventDetails(builder *strings.Builder) {
	builder.WriteString(fmt.Sprintf("DTSTAMP:%s", fmt.Sprintf("%sZ", timeToICal(time.Now().UTC()))) + lineBreak)
	builder.WriteString("SUMMARY:" + e.Title + lineBreak)
	builder.WriteString("LOCATION:" + e.Location + lineBreak)

	if e.Description != "" {
		description := cleanDescription(e.Description)
		builder.WriteString("DESCRIPTION:" + description + lineBreak)
	}

}

func (e *Event) HasRecurrences() bool {
	return len(e.Recurrences) > 0
}
func (e *Event) ConflictsWith(other *Event) (bool, time.Weekday) {
	//Both are recurring events
	if e.HasRecurrences() && other.HasRecurrences() {
		for _, rec := range e.Recurrences {
			for _, otherRec := range other.Recurrences {
				if conflicts, day := rec.ConflictsWith(otherRec); conflicts {
					return true, day
				}
			}
		}
	}

	// Single event vs single event
	if !e.HasRecurrences() && !other.HasRecurrences() {
		if e.StartDate.Before(other.EndDate) && e.EndDate.After(other.StartDate) {
			return true, e.StartDate.Weekday()
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
					return true, singleDay
				}
			}
		}
	}

	return false, 0
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
	if len(desc) > 63 {
		description = desc[:60] + "..."
	} else {
		description = desc
	}
	return strings.ReplaceAll(description, "\n", " ")
}

func timeToICal(t time.Time) string {
	return t.Format(iCalTimeLayout)
}

func stripDay(t time.Time) time.Time {
	return time.Date(0, 1, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}
