package ical

import (
	"fmt"
	"strings"
	"time"
)

type Recurrences struct {
	Frequency  Frequency
	Day        time.Weekday
	StartTime  time.Time
	EndTime    time.Time
	Exceptions []time.Time
}

func (r *Recurrences) Generate(startDate, endDate time.Time, timeZone TimeZone) (string, error) {
	if !r.Valid() {
		return "", ErrInvalidRecurrence
	}

	startTime, err := findStartDate(startDate, r.Day, r.StartTime)
	if err != nil {
		return "", err
	}
	endTime, err := findEndDate(endDate, r.Day, r.EndTime)
	if err != nil {
		return "", err
	}
	if endTime.Before(startTime) {
		return "", ErrEndTimeBeforeStartTime(r.EndTime.String(), r.StartTime.String())
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("DTSTART;TZID=%s:%s", timeZone.iCal(), timeToICal(startTime)))
	builder.WriteString(lineBreak)
	builder.WriteString(fmt.Sprintf("DTEND;TZID=%s:%s", timeZone.iCal(), timeToICal(startTime.Add(r.EndTime.Sub(r.StartTime)))))
	builder.WriteString(lineBreak)
	builder.WriteString("RRULE:")
	builder.WriteString(fmt.Sprintf("FREQ=%s;", r.Frequency))
	builder.WriteString(fmt.Sprintf("BYDAY=%s;", weekdayToICal(r.Day)))
	builder.WriteString(fmt.Sprintf("UNTIL=%s;", fmt.Sprintf("%sZ", timeToICal(endTime.UTC()))))
	if len(r.Exceptions) > 0 {
		for _, ex := range r.Exceptions {
			builder.WriteString(lineBreak)
			builder.WriteString(fmt.Sprintf("EXDATE;TZID=%s:%s", timeZone.iCal(), timeToICal(ex)))
		}
	}

	return builder.String(), nil
}

func (r *Recurrences) Valid() bool {
	if !r.Frequency.Valid() {
		return false
	}
	if !validWeekday(r.Day) {
		return false
	}
	if r.EndTime.Before(r.StartTime) || r.EndTime.Equal(r.StartTime) {
		return false
	}

	return true
}
func (r *Recurrences) ConflictsWith(other Recurrences) (bool, time.Weekday) {
	if r.Day != other.Day {
		return false, 0
	}
	if stripDay(r.StartTime).Before(stripDay(other.EndTime)) && stripDay(other.StartTime).Before(stripDay(r.EndTime)) {
		return true, r.Day
	}
	return false, 0
}

func validWeekday(d time.Weekday) bool {
	switch d {
	case time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday:
		return true
	default:
		return false
	}
}

func DayOfWeekFromString(day string) (time.Weekday, error) {
	switch strings.ToLower(day) {
	case "monday", "mo", "mon":
		return time.Monday, nil
	case "tuesday", "tu", "tue":
		return time.Tuesday, nil
	case "wednesday", "we", "wed":
		return time.Wednesday, nil
	case "thursday", "th", "thu":
		return time.Thursday, nil
	case "friday", "fr", "fri":
		return time.Friday, nil
	case "saturday", "sa", "sat":
		return time.Saturday, nil
	case "sunday", "su", "sun":
		return time.Sunday, nil
	default:
		return 0, ErrInvalidDayOfWeek
	}
}

func weekdayToICal(d time.Weekday) string {
	switch d {
	case time.Monday:
		return "MO"
	case time.Tuesday:
		return "TU"
	case time.Wednesday:
		return "WE"
	case time.Thursday:
		return "TH"
	case time.Friday:
		return "FR"
	case time.Saturday:
		return "SA"
	case time.Sunday:
		return "SU"
	default:
		return ""
	}
}

func findStartDate(startDate time.Time, day time.Weekday, startTime time.Time) (time.Time, error) {

	if !validWeekday(day) {
		return time.Time{}, ErrInvalidDayOfWeek
	}

	currentDate := startDate.In(startDate.Location())

	for currentDate.Weekday() != day {
		currentDate = currentDate.AddDate(0, 0, 1)

		if currentDate.Sub(startDate) > (7 * 24 * time.Hour) {
			return time.Time{}, fmt.Errorf("recurrences/findStartDate: failed to find target day within one week")
		}
	}

	firstOccurrence := time.Date(
		currentDate.Year(),
		currentDate.Month(),
		currentDate.Day(),
		startTime.Hour(),
		startTime.Minute(),
		startTime.Second(),
		0,
		currentDate.Location(),
	)

	return firstOccurrence, nil
}

func findEndDate(endDate time.Time, day time.Weekday, endTime time.Time) (time.Time, error) {

	if !validWeekday(day) {
		return time.Time{}, ErrInvalidDayOfWeek
	}

	currentDate := endDate.In(endDate.Location())

	for currentDate.Weekday() != day {
		currentDate = currentDate.AddDate(0, 0, -1)

		if endDate.Sub(currentDate) > (7 * 24 * time.Hour) {
			return time.Time{}, fmt.Errorf("recurrences/findEndDate: failed to find target day within one week")
		}
	}

	lastOccurrence := time.Date(
		currentDate.Year(),
		currentDate.Month(),
		currentDate.Day(),
		endTime.Hour(),
		endTime.Minute(),
		endTime.Second(),
		0,
		currentDate.Location(),
	)

	return lastOccurrence, nil
}
