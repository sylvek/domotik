package compute

type Output struct {
	WattPerHourForLastMinute float64
	WattPerHourForThisHour   int64
	WattConsumedToday        int64
	EuroSpentToday           float64
	RatioLowTariffToday      float64
}

type Rule interface {
	Tick() chan Output
	Stop()
	SetLowTariffState(state bool)
	SetIndice(indice int64)
}
