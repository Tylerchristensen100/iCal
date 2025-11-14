package ical

import (
	"strings"
	"time"
)

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
