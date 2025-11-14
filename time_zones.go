package ical

import timezones "github.com/Tylerchristensen100/iCal/timezones"

// TimeZone represents the time zone for an event.
// It includes all valid time zones defined by the IANA Time Zone Database.
// All TimeZones are copied from https://github.com/Tylerchristensen100/iCal_VTIMEZONE
//
// iCalendar VTIMEZONE component
type TimeZone timezones.TZID

// Return the TimeZone Definition as specified by the ICal standard.
//
// https://icalendar.org/iCalendar-RFC-5545/3-2-19-time-zone-identifier.html
//
// Timezone information is provided by https://github.com/Tylerchristensen100/iCal_VTIMEZONE
func (tz *TimeZone) iCal() (string, bool) {
	if !tz.valid() {
		return string(timezones.UTC), false
	}

	return timezones.Get(timezones.TZID(*tz))
}

// Return the ID of the TimeZone
//
// E.g., "America/New_York"
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
