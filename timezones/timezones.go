package timezones

import _ "embed"

//go:embed us_central.ics
var USCentral []byte

//go:embed us_eastern.ics
var USEastern []byte

//go:embed us_mountain.ics
var USMountain []byte

//go:embed us_pacific.ics
var USPacific []byte

//go:embed us_alaska.ics
var USAlaska []byte

//go:embed us_arizona.ics
var USArizona []byte

//go:embed us_hawaii.ics
var USHawaii []byte

//go:embed utc.ics
var UTC []byte
