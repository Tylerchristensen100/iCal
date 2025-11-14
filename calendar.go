package ical

import (
	"os"
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
	Journals    []Journal
	Todos       []Todo
}

func (c *Calendar) AddEvent(e Event) error {
	if !e.Valid() {
		return ErrInvalidEvent
	}
	c.Events = append(c.Events, e)
	return nil
}

func (c *Calendar) AddJournal(j Journal) error {
	if !j.valid() {
		return ErrInvalidJournal
	}
	c.Journals = append(c.Journals, j)
	return nil
}

func (c *Calendar) AddTodo(t Todo) error {
	if !t.valid() {
		return ErrInvalidTodo
	}
	c.Todos = append(c.Todos, t)
	return nil
}

// Saves the calendar to a file at the specified path.
func (c *Calendar) Save(path string) error {
	icalData, err := c.Generate()
	if err != nil {
		return err
	}

	if !strings.Contains(path, ".ics") {
		path += ".ics"
	}

	err = os.WriteFile(path, []byte(icalData), 0644)
	if err != nil {
		return err
	}
	return nil
}

// Generate creates the iCal formatted string for the entire calendar.
func (c *Calendar) Generate() ([]byte, error) {
	if !c.Valid() {
		return nil, ErrInvalidCalendar
	}

	var builder strings.Builder
	builder.WriteString("BEGIN:VCALENDAR" + lineBreak)
	builder.WriteString("VERSION:2.0" + lineBreak)
	builder.WriteString("PRODID:-//UVU//Class Schedule//EN" + lineBreak)
	builder.WriteString("CALSCALE:GREGORIAN" + lineBreak)
	builder.WriteString("METHOD:PUBLISH" + lineBreak)

	c.generateTimeZones(&builder)

	for _, event := range c.Events {
		event, err := event.Generate()
		if err != nil {
			return nil, err
		}
		builder.WriteString(event)
	}
	for _, journal := range c.Journals {
		err := journal.generate(&builder)
		if err != nil {
			return nil, err
		}
	}

	for _, todo := range c.Todos {
		err := todo.generate(&builder)
		if err != nil {
			return nil, err
		}
	}

	builder.WriteString("END:VCALENDAR" + lineBreak)
	return []byte(builder.String()), nil
}

func (c *Calendar) generateTimeZones(builder *strings.Builder) {
	uniqueTimeZones := make(map[TimeZone]bool)
	for _, event := range c.Events {
		uniqueTimeZones[event.TimeZone] = true
	}
	for timeZone := range uniqueTimeZones {
		data, found := timeZone.iCal()
		if found {
			builder.WriteString(data + lineBreak)
		}
	}
}

func (c *Calendar) Valid() bool {
	if c.Name == "" {
		return false
	}
	for _, event := range c.Events {
		if !event.Valid() {
			return false
		}
	}
	for _, journal := range c.Journals {
		if !journal.valid() {
			return false
		}
	}
	for _, todo := range c.Todos {
		if !todo.valid() {
			return false
		}
	}

	// If there is nothing in the calendar, it's invalid
	if len(c.Events) == 0 && len(c.Journals) == 0 && len(c.Todos) == 0 {
		return false
	}
	return true
}

// ListConflicts returns a list of events that have scheduling conflicts with other events in the calendar.
func (c *Calendar) ListConflicts() []*Event {
	var conflicts []*Event
	c.ResolveConflicts(func(event1, event2 *Event, _ time.Time) {
		conflicts = append(conflicts, event1, event2)
	})
	return conflicts
}

// ResolveConflicts checks for scheduling conflicts between events in the calendar
// and applies the provided resolveFunc to each pair of conflicting events.
//
// This is for interactive conflict resolution where the user can define how to handle conflicts.
func (c *Calendar) ResolveConflicts(resolveFunc ResolveFunc) {
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

// Function to resolve conflicts between two events.
//
// event1 and event2 are the conflicting events,
// conflictingDay is the day on which the conflict occurs.
//
// Used for more interactive conflict resolution.
/*
 func exampleResolveFunc(event1, event2 *Event, conflictingDay time.Time) {
	if event1.HasRecurrences() && !event2.HasRecurrences() {
		event1.AddException(conflictingDay)
	}
 }
*/
type ResolveFunc func(event1, event2 *Event, conflictingDay time.Time)
