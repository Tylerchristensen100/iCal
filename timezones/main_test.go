package timezones

import "testing"

func TestGet(t *testing.T) {
	tests := []struct {
		tzid       TZID
		shouldFind bool
	}{
		{America_New_York, true},
		{America_Los_Angeles, true},
		{UTC, true},
		{TZID("Invalid/Timezone"), false},
	}

	for _, tt := range tests {
		data, found := Get(tt.tzid)
		if found != tt.shouldFind {
			t.Errorf("Get(%q) = found %v; want %v", tt.tzid, found, tt.shouldFind)
		}
		if found && data == "" {
			t.Errorf("Get(%q) returned empty data", tt.tzid)
		}
	}
}

func TestGetFailure(t *testing.T) {
	tzid := TZID("NonExistent/Timezone")
	data, found := Get(tzid)
	if found {
		t.Errorf("Get(%q) = found %v; want false", tzid, found)
	}
	if data != "" {
		t.Errorf("Get(%q) returned data %q; want empty string", tzid, data)
	}
}

func TestValid(t *testing.T) {
	tests := []struct {
		tzid    TZID
		isValid bool
	}{
		{America_Chicago, true},
		{America_Denver, true},
		{TZID("Invalid/Timezone"), false},
	}

	for _, tt := range tests {
		if tt.tzid.Valid() != tt.isValid {
			t.Errorf("Valid(%q) = %v; want %v", tt.tzid, tt.tzid.Valid(), tt.isValid)
		}
	}
}

func TestLoad(t *testing.T) {
	err := load()
	if err != nil {
		t.Errorf("load() returned error: %v", err)
	}
	if len(timezones) == 0 {
		t.Errorf("load() did not populate timezones map")
	}
}
