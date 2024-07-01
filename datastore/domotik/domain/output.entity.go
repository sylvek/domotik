package domain

type Output struct {
	wattPerHourForLastMinute float64
	wattPerHourForThisHour   int64
	wattConsumedToday        int64
	euroSpentToday           float64
	ratioLowTariffToday      float64
}

func NewOutput(dailySumLow int64, dailySumHigh int64, hourlySum int64, hourlyNbIndices int64, consumptionSinceLastTime int64, minutesSinceTheLastIndice float64) *Output {

	HIGH_TARIFF_PRICE := 0.0001963
	LOW_TARIFF_PRICE := 0.0001457

	ratioLowTariffToday := 1.0
	if dailySumHigh > 0 {
		ratioLowTariffToday = float64(dailySumLow) / float64(dailySumHigh+dailySumLow)
	}

	return &Output{
		wattPerHourForLastMinute: float64(consumptionSinceLastTime) * 60 / minutesSinceTheLastIndice,
		wattPerHourForThisHour:   hourlySum * 60 / hourlyNbIndices,
		wattConsumedToday:        dailySumHigh + dailySumLow,
		euroSpentToday:           float64(dailySumHigh)*HIGH_TARIFF_PRICE + float64(dailySumLow)*LOW_TARIFF_PRICE,
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
