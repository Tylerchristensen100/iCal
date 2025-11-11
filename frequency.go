package ical

// Frequency represents the frequency of recurrence for an event.
// It can be daily, weekly, monthly, or yearly.
type Frequency string

const (
	DailyFrequency   Frequency = "DAILY"
	WeeklyFrequency  Frequency = "WEEKLY"
	MonthlyFrequency Frequency = "MONTHLY"
	YearlyFrequency  Frequency = "YEARLY"
)

func (f *Frequency) Valid() bool {
	switch *f {
	case DailyFrequency, WeeklyFrequency, MonthlyFrequency, YearlyFrequency:
		return true
	default:
		return false
	}
}
