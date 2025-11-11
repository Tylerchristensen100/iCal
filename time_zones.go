package ical

type TimeZone string

const (
	MountainTimeZone TimeZone = "America/Denver"
	EasternTimeZone  TimeZone = "America/New_York"
	CentralTimeZone  TimeZone = "America/Chicago"
	PacificTimeZone  TimeZone = "America/Los_Angeles"
	AlaskaTimeZone   TimeZone = "America/Anchorage"
	HawaiiTimeZone   TimeZone = "America/Honolulu"
	ArizonaTimeZone  TimeZone = "America/Phoenix"
	UTC              TimeZone = "UTC"
)

func (tz *TimeZone) iCal() string {
	// Return the TimeZone Definition as specified by the ICal standard.
	// https://icalendar.org/iCalendar-RFC-5545/3-2-19-time-zone-identifier.html

	if !tz.Valid() {
		return string(UTC)
	}

	return string(*tz)
}

func (tz *TimeZone) Valid() bool {
	switch *tz {
	case MountainTimeZone, EasternTimeZone, CentralTimeZone, PacificTimeZone,
		AlaskaTimeZone, HawaiiTimeZone, ArizonaTimeZone, UTC:
		return true
	default:
		return false
	}
}
