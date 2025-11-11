package ical

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
