package ical

import (
	"fmt"
	"strings"
	"time"
)

// Recurrences represents the recurrence rules for an event.
type Recurrences struct {

	// REQUIRED: Frequency of the recurrence
	Frequency Frequency

	// REQUIRED: Day of the week for the recurrence
	Day time.Weekday

	// REQUIRED: Start and end time for each occurrence
	StartTime time.Time

	// REQUIRED: End time for each occurrence
	EndTime time.Time

	// OPTIONAL: List of exception dates for the recurrence
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
	builder.WriteString(fmt.Sprintf("DTSTART;TZID=%s:%s", timeZone.ID(), timeToICal(startTime)))
	builder.WriteString(lineBreak)
	builder.WriteString(fmt.Sprintf("DTEND;TZID=%s:%s", timeZone.ID(), timeToICal(startTime.Add(r.EndTime.Sub(r.StartTime)))))
	builder.WriteString(lineBreak)
	builder.WriteString(generateRRULE(r.Frequency, r.Day, endTime))
	if len(r.Exceptions) > 0 {
		for _, ex := range r.Exceptions {
			builder.WriteString(lineBreak)
			builder.WriteString(fmt.Sprintf("EXDATE;TZID=%s:%s", timeZone.ID(), timeToICal(ex)))
		}
	}

	return builder.String(), nil
}

func generateRRULE(f Frequency, d time.Weekday, endTime time.Time) string {
	return fmt.Sprintf("RRULE:FREQ=%s;BYDAY=%s;UNTIL=%s;"+lineBreak, f, weekdayToICal(d), fmt.Sprintf("%sZ", timeToICal(endTime.UTC())))
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

func (r *Recurrences) uid() string {
	return fmt.Sprintf("%s-%s-%s@iCal.go", r.Frequency,
		r.StartTime.UTC().Format("15_04"), r.EndTime.UTC().Format("15_04"))
}
func (r *Recurrences) ConflictsWith(other Recurrences) (bool, time.Time) {
	if r.Day != other.Day {
		return false, time.Time{}
	}
	if stripDay(r.StartTime).Before(stripDay(other.EndTime)) && stripDay(other.StartTime).Before(stripDay(r.EndTime)) {
		return true, time.Date(0, 1, 1, r.StartTime.Hour(), r.StartTime.Minute(), r.StartTime.Second(), 0, time.UTC)
	}

	return false, time.Time{}
}

func (r *Recurrences) Occurrences(startDate time.Time, endDate time.Time) []time.Time {
	var occurrences []time.Time
	startTime := stripDay(r.StartTime)
	startDate, err := findStartDate(startDate, r.Day, r.StartTime)
	if err != nil {
		return occurrences
	}

	current := startDate
	for current.Before(endDate) || current.Equal(endDate) {
		occurrence := time.Date(
			current.Year(),
			current.Month(),
			current.Day(),
			startTime.Hour(),
			startTime.Minute(),
			startTime.Second(),
			0,
			current.Location(),
		)
		occurrences = append(occurrences, occurrence)

		switch r.Frequency {
		case DailyFrequency:
			current = current.AddDate(0, 0, 1)
		case WeeklyFrequency:
			current = current.AddDate(0, 0, 7)
		case MonthlyFrequency:
			current = current.AddDate(0, 1, 0)
		case YearlyFrequency:
			current = current.AddDate(1, 0, 0)
		}
	}
	return occurrences
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

func validWeekday(d time.Weekday) bool {
	switch d {
	case time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday:
		return true
	default:
		return false
	}
}
