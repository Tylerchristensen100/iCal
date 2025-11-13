package timezones

// Copied and adapted from  https://github.com/Tylerchristensen100/iCal_VTIMEZONE/blob/main/examples/golang/main.go
// REPO: https://github.com/Tylerchristensen100/iCal_VTIMEZONE

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/json"
	"fmt"
)

// from `/tz/timezones.json`
//
//go:embed timezones.json.gz
var compressedData []byte

var timezones map[string]string

func init() {
	err := load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load timezones: %v", err))
	}
}

func load() error {
	reader, err := gzip.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer reader.Close()

	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&timezones); err != nil {
		return fmt.Errorf("failed to decode JSON data: %w", err)
	}
	return nil
}

func Get(tzid TZID) (string, bool) {
	data, exists := timezones[string(tzid)]
	if !exists {
		return "", false
	}

	return data, exists
}

func (tz *TZID) Valid() bool {
	_, exists := timezones[string(*tz)]
	return exists
}
