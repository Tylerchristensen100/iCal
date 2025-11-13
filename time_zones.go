package ical

import timezones "github.com/Tylerchristensen100/iCal/timezones"

// TimeZone represents the time zone for an event.
// It includes common time zones used in the United States and UTC.
//
// iCalendar VTIMEZONE component
type TimeZone timezones.TZID

// Return the TimeZone Definition as specified by the ICal standard.
// https://icalendar.org/iCalendar-RFC-5545/3-2-19-time-zone-identifier.html
func (tz *TimeZone) iCal() (string, bool) {
	if !tz.valid() {
		return string(timezones.UTC), false
	}

	return timezones.Get(timezones.TZID(*tz))
}

func (tz *TimeZone) ID() string {
	return string(*tz)
}

func (tz *TimeZone) valid() bool {
	if tz == nil {
		return false
	}
	zone := timezones.TZID(*tz)
	return zone.Valid()
}
