# iCal

iCal is a Go library for generating iCalendar (.ics) files with support for recurring events. It provides an easy-to-use API for creating events, defining recurrence rules, and exporting calendars in the standard iCalendar format.

[![Go Reference](https://pkg.go.dev/badge/github.com/Tylerchristensen100/iCal.svg)](https://pkg.go.dev/github.com/Tylerchristensen100/iCal)
[![codecov](https://codecov.io/github/Tylerchristensen100/iCal/graph/badge.svg?token=CXIYJTVZV1)](https://codecov.io/github/Tylerchristensen100/iCal)

## Features

- Create events with start and end times, summaries, descriptions, and locations.
- Define recurrence rules for events (daily, weekly, monthly, yearly).
- Handle 1 time exceptions for recurring events.
- Support for To-Do and Journal components.
- Set reminders for events with various actions (display, email, audio).
- Support for all major time zones via the iCal_VTIMEZONE library.
- Export calendars to .ics files compatible with popular calendar applications.

## Installation

To install the iCal library, use the following command:

```bash
go get github.com/Tylerchristensen100/iCal
```

### Creating a Calendar

Here is a simple example of how to create a calendar with a recurring event:

```go
package main

import (
	"os"
	"strings"
	"testing"
	"time"

	ical "github.com/Tylerchristensen100/iCal"
)

func main() {
    calendar := ical.Create("Team Calendar", "Coordinate team meetings and events")

    event := ical.Event{
        Summary:     "Weekly Meeting",
        Description: "Team sync-up",
        Location:    "Conference Room",
        StartTime:   time.Date(2025, time.July, 1, 9, 0, 0, 0, time.UTC),
        EndTime:     time.Date(2026, time.September, 1, 10, 0, 0, 0, time.UTC),
        Recurrences: ical.Recurrences{
            Frequency: ical.WeeklyFrequency,
            Day:       time.Thursday,
            StartTime: time.Date(2025, time.July, 1, 9, 0, 0, 0, time.UTC),
            EndTime:   time.Date(2025, time.July, 1, 10, 0, 0, 0, 0, time.UTC),
            Exceptions: []time.Time{
                time.Date(2025, time.November, 27, 9, 0, 0, 0, time.UTC), // Skip Thanksgiving
        },
    }
    }

    thanksgiving := ical.Event{
        Summary:     "Thanksgiving",
        Description: "Time Off",
        Location:    "Anywhere but the office",
        StartTime:   time.Date(2025, time.November, 27, 0, 0, 0, 0, time.UTC),
        EndTime:     time.Date(2025, time.November, 27, 23, 59, 59, 0, time.UTC),
    }

    calendar.AddEvent(event)
    calendar.AddEvent(thanksgiving)

    icsData, err := calendar.Generate()
	if err != nil {
		panic(err)
	}

    err = os.WriteFile("team_calendar.ics", []byte(icsData), 0644)
	if err != nil {
		panic(err)
	}
}

```



## Notes
- TimeZones are provided via the [iCal_VTIMEZONE](https://github.com/Tylerchristensen100/iCal_VTIMEZONE) library.
- Timezones are embedded directly into a map within this library for ease of use.  This means there is a 104kb increase in binary size.  The total size of the built library is approximately 550Kb.



## Resources
 - [iCalendar Specification (RFC 5545)](https://icalendar.org/RFC-Specifications/iCalendar-RFC-5545/)