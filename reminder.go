package ical

import (
	"fmt"
	"strings"
	"time"
)

// iCalendar VALARM component
//
// https://icalendar.org/iCalendar-RFC-5545/3-6-6-alarm-component.html
type Reminder struct {

	// REQUIRED: Description of the reminder
	Description string

	// REQUIRED: Action to be taken when the reminder is triggered
	//
	// Possible Values: DisplayReminderAction, EmailReminderAction, AudioReminderAction
	//
	// **EmailReminderAction** requires SUMMARY & ATTENDEES properties
	Action ReminderAction

	// REQUIRED: Time before the event when the reminder should trigger
	Trigger time.Duration

	// OPTIONAL: Number of times the reminder should repeat
	Repeat *int

	// OPTIONAL: List of attendees to notify
	//
	// **Only for EMAIL action**
	Attendees []Participant
}

// The Action to be taken when the reminder is triggered
type ReminderAction string

const (
	DisplayReminderAction ReminderAction = "DISPLAY"

	// Requires SUMMARY & ATTENDEES properties
	EmailReminderAction ReminderAction = "EMAIL"

	AudioReminderAction ReminderAction = "AUDIO"
)

// Pretty print a user friendly representation of the reminder
func (r *Reminder) String() string {
	builder := strings.Builder{}
	builder.WriteString("Reminder:\n")
	builder.WriteString("  Description: " + r.Description + "\n")
	builder.WriteString("  Action: " + string(r.Action) + "\n")
	builder.WriteString("  Trigger: " + r.Trigger.String() + "\n")
	if r.Repeat != nil {
		builder.WriteString("  Repeats: " + fmt.Sprintf("%d", *r.Repeat) + " times\n")
	}

	if r.Action == EmailReminderAction && len(r.Attendees) > 0 {
		builder.WriteString("  Attendees:\n")
		for _, attendee := range r.Attendees {
			builder.WriteString("    - name: " + attendee.Name + ", email: " + attendee.Email + "\n")
		}
	}
	return builder.String()
}

// Generate iCalendar VALARM component
func (r *Reminder) generate(builder *strings.Builder) error {
	if !r.valid() {
		return ErrInvalidReminder
	}
	builder.WriteString("BEGIN:VALARM\r\n")
	builder.WriteString("ACTION:" + string(r.Action) + "\r\n")
	builder.WriteString("DESCRIPTION:" + cleanDescription(r.Description) + "\r\n")
	builder.WriteString("TRIGGER:" + formatDurationAsTrigger(r.Trigger) + "\r\n")
	if r.Repeat != nil {
		builder.WriteString("REPEAT:" + fmt.Sprintf("%d", *r.Repeat) + "\r\n")
		// Assuming a fixed DURATION of 15 minutes for each repeat for simplicity
		builder.WriteString("DURATION:PT15M\r\n")
	}
	if r.Action == EmailReminderAction && len(r.Attendees) > 0 {
		for _, attendee := range r.Attendees {
			builder.WriteString("ATTENDEE;CN=" + attendee.Name + ":MAILTO:" + attendee.Email + "\r\n")
		}
	}
	builder.WriteString("END:VALARM\r\n")
	return nil
}

func (r *Reminder) addTrigger(duration time.Duration) {
	r.Trigger = duration
}

func (r *Reminder) valid() bool {
	if r.Description == "" {
		return false
	}
	if r.Action != DisplayReminderAction && r.Action != EmailReminderAction && r.Action != AudioReminderAction {
		return false
	}

	if r.Action == EmailReminderAction && len(r.Attendees) == 0 {
		return false
	}
	if r.Action != EmailReminderAction && len(r.Attendees) > 0 {
		return false
	}

	if r.Repeat != nil && *r.Repeat < 0 {
		return false
	}

	if r.Trigger == 0 {
		return false
	}

	return true
}

// Format time.Duration as iCal TRIGGER value
func formatDurationAsTrigger(d time.Duration) string {
	// iCal duration format is using ISO 8601 Durations
	// P(n)Y(n)M(n)DT(n)H(n)M(n)S
	isNegative := d < 0
	if isNegative {
		d = -d
	}

	totalSeconds := int(d.Seconds())
	seconds := totalSeconds % 60
	minutes := (totalSeconds / 60) % 60
	hours := (totalSeconds / (60 * 60))
	days := (totalSeconds / (60 * 60)) / 24

	if days > 0 {
		hours = hours % 24
	}

	var builder strings.Builder
	if isNegative {
		builder.WriteString("-")
	}
	builder.WriteString("P")

	if days > 0 {
		builder.WriteString(fmt.Sprintf("%dD", days))
	}

	hasTime := hours > 0 || minutes > 0 || seconds > 0
	if hasTime {
		builder.WriteString("T")
	}

	if hours > 0 {
		builder.WriteString(fmt.Sprintf("%dH", hours))
	}
	if minutes > 0 {
		builder.WriteString(fmt.Sprintf("%dM", minutes))
	}
	if seconds > 0 || !hasTime {
		builder.WriteString(fmt.Sprintf("%dS", seconds))
	}

	return builder.String()
}
