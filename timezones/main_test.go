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
