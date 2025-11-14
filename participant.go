package ical

import (
	"fmt"
	"strings"
)

type Participant struct {
	Name  string
	Email string
}

func (p *Participant) generate(builder *strings.Builder) error {
	if !p.valid() {
		return ErrInvalidEmail
	}
	builder.WriteString(fmt.Sprintf("ATTENDEE;CUTYPE=INDIVIDUAL;ROLE=REQ-PARTICIPANT;PARTSTAT=NEEDS-ACTION;RSVP=TRUE;CN=%s;X-NUM-GUESTS=0:mailto:%s", p.Name, p.Email) + lineBreak)
	return nil
}

func (p *Participant) generateOrganizer(builder *strings.Builder) error {
	if !p.valid() {
		return ErrInvalidEmail
	}
	builder.WriteString(fmt.Sprintf("ORGANIZER;CN=%s:mailto:%s", p.Name, p.Email) + lineBreak)
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
