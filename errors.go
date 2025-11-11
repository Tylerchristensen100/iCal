package ical

import "fmt"

const (
	errInvalidEventMessage      = "invalid event"
	errInvalidRecurrenceMessage = "invalid recurrence"
	errInvalidDayOfWeekMessage  = "invalid day of week"
)

var (
	ErrInvalidEvent      = fmt.Errorf(errInvalidEventMessage)
	ErrInvalidRecurrence = fmt.Errorf(errInvalidRecurrenceMessage)
	ErrInvalidDayOfWeek  = fmt.Errorf(errInvalidDayOfWeekMessage)
)

func ErrEndTimeBeforeStartTime(endTime, startTime string) error {
	return fmt.Errorf("end time '%s' is before start time '%s'", endTime, startTime)
}
