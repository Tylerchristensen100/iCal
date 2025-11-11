package ical

import (
	"strings"
	"time"
)

const (
	lineBreak      = "\r\n"
	iCalTimeLayout = "20060102T150405"
)

type Calendar struct {
	Name        string
	Description string
	Events      []Event
}

func (c *Calendar) AddEvent(e Event) error {
	if !e.Valid() {
		return ErrInvalidEvent
	}
	c.Events = append(c.Events, e)
	return nil
}

// Generate creates the iCal formatted string for the entire calendar.
func (c *Calendar) Generate() (string, error) {
	var builder strings.Builder
	builder.WriteString("BEGIN:VCALENDAR" + lineBreak)
	builder.WriteString("VERSION:2.0" + lineBreak)
	builder.WriteString("PRODID:-//UVU//Class Schedule//EN" + lineBreak)
	builder.WriteString("CALSCALE:GREGORIAN" + lineBreak)
	builder.WriteString("METHOD:PUBLISH" + lineBreak)

	for _, event := range c.Events {
		event, err := event.Generate()
		if err != nil {
			return "", err
		}
		builder.WriteString(event)
	}
	builder.WriteString("END:VCALENDAR" + lineBreak)
	return builder.String(), nil
}

// ListConflicts returns a list of events that have scheduling conflicts with other events in the calendar.
func (c *Calendar) ListConflicts() []*Event {
	var conflicts []*Event
	c.ResolveConflicts(func(event1, event2 *Event, _ time.Weekday) {
		conflicts = append(conflicts, event1, event2)
	})
	return conflicts
}

// ResolveConflicts checks for scheduling conflicts between events in the calendar
// and applies the provided resolveFunc to each pair of conflicting events.
//
// This is for interactive conflict resolution where the user can define how to handle conflicts.
func (c *Calendar) ResolveConflicts(resolveFunc func(event1, event2 *Event, conflictingDay time.Weekday)) {
	for i, event := range c.Events {
		for j, otherEvent := range c.Events {
			// Skip self-comparison
			if i == j {
				continue
			}
			if conflicts, day := event.ConflictsWith(&otherEvent); conflicts {
				resolveFunc(&c.Events[i], &c.Events[j], day)
			}
		}
	}
}
