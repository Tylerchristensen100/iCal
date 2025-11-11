package ical

import "fmt"

const (
	errInvalidEventMessage      = "invalid event"
	errInvalidRecurrenceMessage = "invalid recurrence"
	errInvalidDayOfWeekMessage  = "invalid day of week"
)

var (
	// ErrInvalidEvent is returned when an event is not valid.
	ErrInvalidEvent = fmt.Errorf(errInvalidEventMessage)

	// ErrInvalidRecurrence is returned when a recurrence rule is not valid.
	ErrInvalidRecurrence = fmt.Errorf(errInvalidRecurrenceMessage)
	// ErrInvalidDayOfWeek is returned when a day of the week is not valid.
	ErrInvalidDayOfWeek = fmt.Errorf(errInvalidDayOfWeekMessage)
)

// ErrEndTimeBeforeStartTime is returned when the end time is before the start time.
func ErrEndTimeBeforeStartTime(endTime, startTime string) error {
	return fmt.Errorf("end time '%s' is before start time '%s'", endTime, startTime)
}
