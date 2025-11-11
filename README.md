# iCal

iCal is a Go library for generating iCalendar (.ics) files with support for recurring events. It provides an easy-to-use API for creating events, defining recurrence rules, and exporting calendars in the standard iCalendar format.

## Features

- Create events with start and end times, summaries, descriptions, and locations.
- Define recurrence rules for events (daily, weekly, monthly, yearly).
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
