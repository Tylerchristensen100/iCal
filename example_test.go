package ical

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Tylerchristensen100/iCal/timezones"
)

func ExampleCalendar_AddEvent() {
	cal := Create("Example Calendar", "An example calendar")
	event := Event{
		Title:     "Meeting with Bob",
		StartDate: time.Date(2024, 7, 1, 10, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 7, 1, 11, 0, 0, 0, time.UTC),
		Organizer: &Participant{Name: "Organizer", Email: "example@github.com"},
		TimeZone:  TimeZone(timezones.UTC),
	}
	err := cal.AddEvent(event)
	if err != nil {
		panic(err)
	}

	output, err := cal.Generate()
	if err != nil {
		panic(err)
	}

	// Remove the generated attributes for consistent output
	re := regexp.MustCompile(`DTSTAMP:\d{8}T\d{6}Z\n?`)
	regUid := regexp.MustCompile(`(UID:[^\n]*?)-\d+@iCal\.go`)
	// Generated value at generation time, replace with fixed value
	validOutput := re.ReplaceAllString(string(output), "DTSTAMP:20251114T212240Z")
	// Remove the random unique int in the UID for testing
	validOutput = regUid.ReplaceAllString(validOutput, "$1-<unique>@iCal.go")
	// Normalize line endings for consistent output across platforms
	validOutput = strings.ReplaceAll(validOutput, "\r\n", "\n")
	// END the removal of generated attributes

	fmt.Println(validOutput)
	// Output:
	// BEGIN:VCALENDAR
	// VERSION:2.0
	// PRODID:-//TylerChristensen100//iCal_Generator//EN
	// CALSCALE:GREGORIAN
	// METHOD:PUBLISH
	// BEGIN:VTIMEZONE
	// TZID:UTC
	// COMMENT:This timezone only works from 1970-01-01 to 2038-01-01.
	// BEGIN:STANDARD
	// DTSTART:19700101T000000
	// TZNAME:UTC
	// TZOFFSETFROM:+0000
	// TZOFFSETTO:+0000
	// END:STANDARD
	// END:VTIMEZONE
	// BEGIN:VEVENT
	// UID:Meeting_with_Bob-Monday-Monday-<unique>@iCal.go
	// DTSTART;TZID=UTC:20240701T100000
	// DTEND;TZID=UTC:20240701T110000
	// DTSTAMP:20251114T212240Z
	// SUMMARY:Meeting with Bob
	// ORGANIZER;CN=Organizer:mailto:example@github.com
	// END:VEVENT
	// END:VCALENDAR
}

func ExampleCalendar_AddTodo() {
	cal := Create("Example Calendar", "An example calendar")
	todo := Todo{
		Summary:     "Finish Report",
		Description: "Complete the quarterly report.",
		Status:      InProcessStatus,
		Organizer:   Participant{Name: "Manager", Email: "manager@example.com"},
	}
	err := cal.AddTodo(todo)
	if err != nil {
		panic(err)
	}

	output, err := cal.Generate()
	if err != nil {
		panic(err)
	}

	// Remove the generated attributes for consistent output
	regDTSTAMP := regexp.MustCompile(`DTSTAMP:\d{8}T\d{6}\n?`)
	regUid := regexp.MustCompile(`(UID:[^\n]*?)-\d+@iCal\.go`)
	// Generated value at generation time, replace with fixed value
	validOutput := regDTSTAMP.ReplaceAllString(string(output), "DTSTAMP:20251114T212240Z")
	validOutput = regUid.ReplaceAllString(validOutput, "$1-<unique>@iCal.go")
	// Normalize line endings for consistent output across platforms
	validOutput = strings.ReplaceAll(validOutput, "\r\n", "\n")
	// END the removal of generated attributes

	fmt.Println(validOutput)
	// Output:
	// BEGIN:VCALENDAR
	// VERSION:2.0
	// PRODID:-//TylerChristensen100//iCal_Generator//EN
	// CALSCALE:GREGORIAN
	// METHOD:PUBLISH
	// BEGIN:VTODO
	// UID:Finish_Report-<unique>@iCal.go
	// DTSTAMP:20251114T212240Z
	// SUMMARY:Finish Report
	// STATUS:IN-PROCESS
	// DESCRIPTION:Complete the quarterly report.
	// ORGANIZER;CN=Manager:mailto:manager@example.com
	// END:VTODO
	// END:VCALENDAR
}
