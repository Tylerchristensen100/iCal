package ical

import "fmt"

const (
	errInvalidEventMessage      = "invalid event"
	errInvalidRecurrenceMessage = "invalid recurrence"
	errInvalidDayOfWeekMessage  = "invalid day of week"
	errNoConflictFoundMessage   = "no conflict found for the specified date"
	errNoRecurrenceFoundMessage = "no recurrence found for the specified day"
)

var (
	// ErrInvalidEvent is returned when an event is not valid.
	ErrInvalidEvent = fmt.Errorf(errInvalidEventMessage)

	// ErrInvalidRecurrence is returned when a recurrence rule is not valid.
	ErrInvalidRecurrence = fmt.Errorf(errInvalidRecurrenceMessage)

	// ErrInvalidDayOfWeek is returned when a day of the week is not valid.
	ErrInvalidDayOfWeek = fmt.Errorf(errInvalidDayOfWeekMessage)

	// ErrNoConflictFound is returned when no conflict is found for the specified date.
	ErrNoConflictFound = fmt.Errorf(errNoConflictFoundMessage)

	// ErrNoRecurrenceFound is returned when no recurrence is found for the specified day.
	ErrNoRecurrenceFound = fmt.Errorf(errNoRecurrenceFoundMessage)
)

// ErrEndTimeBeforeStartTime is returned when the end time is before the start time.
func ErrEndTimeBeforeStartTime(endTime, startTime string) error {
	return fmt.Errorf("end time '%s' is before start time '%s'", endTime, startTime)
}
