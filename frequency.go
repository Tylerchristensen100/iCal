package ical

// Frequency represents the frequency of recurrence for an event.
// It can be daily, weekly, monthly, or yearly.
type Frequency string

const (
	DailyFrequency    Frequency = "DAILY"
	WeeklyFrequency   Frequency = "WEEKLY"
	BiWeeklyFrequency Frequency = "BIWEEKLY"
	MonthlyFrequency  Frequency = "MONTHLY"
	YearlyFrequency   Frequency = "YEARLY"
)

func (f *Frequency) Valid() bool {
	switch *f {
	case DailyFrequency, WeeklyFrequency, BiWeeklyFrequency, MonthlyFrequency, YearlyFrequency:
		return true
	default:
		return false
	}
}
