package ical

const (
	// Ical Attribute Constants
	prodID        = "PRODID:-//TylerChristensen100//iCal_Generator//EN"
	uid           = "UID:"
	startDateTime = "DTSTART:"
	endDateTime   = "DTEND:"
	urlAttribute  = "URL:"
	timestamp     = "DTSTAMP:"
	organizer     = "ORGANIZER;CN="
	mailto        = ":mailto:"
	description   = "DESCRIPTION:"
	lineBreak     = "\r\n"

	// CALENDAR
	beginCalendar = "BEGIN:VCALENDAR"
	endCalendar   = "END:VCALENDAR"
	iCalVersion   = "VERSION:2.0"
	scale         = "CALSCALE:GREGORIAN"
	method        = "METHOD:PUBLISH"
	attendee      = "ATTENDEE;CUTYPE=INDIVIDUAL;ROLE=REQ-PARTICIPANT;PARTSTAT=NEEDS-ACTION;RSVP=TRUE;CN=%s;X-NUM-GUESTS=0:mailto:%s"

	// EVENT
	beginEvent = "BEGIN:VEVENT"
	endEvent   = "END:VEVENT"
	eventStart = "DTSTART;"
	eventEnd   = "DTEND;"

	//TODO
	beginTodo       = "BEGIN:VTODO"
	endTodo         = "END:VTODO"
	due             = "DUE:"
	completed       = "COMPLETED:"
	priority        = "PRIORITY:"
	percentComplete = "PERCENT-COMPLETE:"

	//ALARM
	beginAlarm = "BEGIN:VALARM"
	endAlarm   = "END:VALARM"
	action     = "ACTION:"
	trigger    = "TRIGGER:"
	repeat     = "REPEAT:"
	duration   = "DURATION:"

	// JOURNAL
	beginJournal = "BEGIN:VJOURNAL"
	endJournal   = "END:VJOURNAL"
	status       = "STATUS:"
	summary      = "SUMMARY:"

	// TimeZone
	timezoneId = "TZID="

	// Recurrence
	rRule = "RRULE:"
	freq  = "FREQ=%s;"
	byDay = "BYDAY=%s;"
	until = "UNTIL=%sZ;"
)
