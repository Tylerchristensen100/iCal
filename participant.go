package ical

import (
	"fmt"
	"strings"
)

type Participant struct {
	// REQUIRED: Name of the participant
	Name string

	// REQUIRED: Email of the participant
	Email string
}

func (p *Participant) generate(builder *strings.Builder) error {
	if !p.valid() {
		return ErrInvalidEmail
	}
	builder.WriteString(fmt.Sprintf(attendee, p.Name, p.Email) + lineBreak)
	return nil
}

func (p *Participant) generateOrganizer(builder *strings.Builder) error {
	if !p.valid() {
		return ErrInvalidEmail
	}
	builder.WriteString(organizer + p.Name + mailto + p.Email + lineBreak)
	return nil
}

func (p *Participant) valid() bool {
	return validateEmail(p.Email) && p.Name != ""
}

func validateEmail(email string) bool {
	if !strings.Contains(email, "@") || strings.HasPrefix(email, "@") || strings.HasSuffix(email, "@") {
		return false
	}
	return true
}
