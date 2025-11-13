package ical

import "github.com/Tylerchristensen100/iCal/timezones"

// TimeZone represents the time zone for an event.
// It includes common time zones used in the United States and UTC.
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

// Return the TimeZone Definition as specified by the ICal standard.
// https://icalendar.org/iCalendar-RFC-5545/3-2-19-time-zone-identifier.html
func (tz *TimeZone) iCal() string {
	if !tz.valid() {
		return string(UTC)
	}

	return string(*tz)
}

func (tz *TimeZone) valid() bool {
	switch *tz {
	case MountainTimeZone, EasternTimeZone, CentralTimeZone, PacificTimeZone,
		AlaskaTimeZone, HawaiiTimeZone, ArizonaTimeZone, UTC:
		return true
	default:
		return false
	}
}

func (tz *TimeZone) ics() []byte {
	switch *tz {
	case CentralTimeZone:
		return timezones.USCentral
	case EasternTimeZone:
		return timezones.USEastern
	case MountainTimeZone:
		return timezones.USMountain
	case PacificTimeZone:
		return timezones.USPacific
	case AlaskaTimeZone:
		return timezones.USAlaska
	case ArizonaTimeZone:
		return timezones.USArizona
	case HawaiiTimeZone:
		return timezones.USHawaii
	case UTC:
		return timezones.UTC
	default:
		return timezones.UTC
	}
}
