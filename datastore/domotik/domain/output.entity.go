package domain

type Output struct {
	wattPerHourForLastMinute float64
	wattPerHourForThisHour   int64
	wattConsumedToday        int64
	euroSpentToday           float64
	ratioLowTariffToday      float64
}

func NewOutput(state State, consumptionSinceLastTime int64, minutesSinceTheLastIndice float64, ratioLowTariffToday float64) *Output {
	return &Output{
		wattPerHourForLastMinute: float64(consumptionSinceLastTime) * 60 / minutesSinceTheLastIndice,
		wattPerHourForThisHour:   state.HourlySum * 60 / state.HourlyNbIndices,
		wattConsumedToday:        state.DailySumHigh + state.DailySumLow,
		euroSpentToday:           float64(state.DailySumHigh)*0.0001963 + float64(state.DailySumLow)*0.0001457,
		ratioLowTariffToday:      ratioLowTariffToday,
	}
}

func (o *Output) GetEuroSpentToday() float64 {
	return o.euroSpentToday
}

func (o *Output) GetRatioLowTariffToday() float64 {
	return o.ratioLowTariffToday
}

func (o *Output) GetWattConsumedToday() float64 {
	return float64(o.wattConsumedToday)
}

func (o *Output) GetWattPerHourForThisHour() float64 {
	return float64(o.wattPerHourForThisHour)
}

func (o *Output) GetWattPerHourForLastMinute() float64 {
	return float64(o.wattPerHourForLastMinute)
}
