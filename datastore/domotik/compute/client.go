package compute

type Output struct {
	WattPerHourForLastMinute int
	WattPerHourForThisHour   int
	WattConsumedToday        int
	EuroSpentToday           float64
	RatioLowPriceToday       float64
}

type Rule interface {
	Tick() chan Output
	Stop()
	SetLowPriceState(state bool)
	SetIndice(indice int)
}
