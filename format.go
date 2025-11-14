package ical

import "strings"

func cleanDescription(desc string) string {
	// Truncate description to 75 characters (including 'DESCRIPTION:' prefix)
	// https://icalendar.org/iCalendar-RFC-5545/3-1-content-lines.html
	//
	// But, from practical experience, some calendar apps handle longer descriptions, we will just add a newline after 75 chars
	var description string
	description = escapeText(desc)
	if len(description) > 63 {
		var builder strings.Builder
		for i, r := range description {
			if i > 0 && i%63 == 0 {
				builder.WriteString("\r\n ")
			}
			builder.WriteRune(r)
		}
		description = builder.String()
	}
	return description
}

func escapeText(text string) string {
	text = strings.ReplaceAll(text, "\\", "\\\\")
	text = strings.ReplaceAll(text, ";", ",")
	text = strings.ReplaceAll(text, "\r", "")
	text = strings.ReplaceAll(text, "\n", " ")
	return text
}
