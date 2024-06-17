package model

type Output struct {
	WattPerHourForLastMinute float64
	WattPerHourForThisHour   int64
	WattConsumedToday        int64
	EuroSpentToday           float64
	RatioLowTariffToday      float64
}
