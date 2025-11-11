package ical

import (
	"fmt"
	"testing"
)

func TestErrEndTimeBeforeStartTime(t *testing.T) {
	endTime := "10:00"
	startTime := "09:00"
	s := ErrEndTimeBeforeStartTime(endTime, startTime).Error()
	expected := fmt.Sprintf("end time '%s' is before start time '%s'", endTime, startTime)
	if s != expected {
		t.Errorf("ErrEndTimeBeforeStartTime() returned '%s', expected '%s'", s, expected)
	}
}
