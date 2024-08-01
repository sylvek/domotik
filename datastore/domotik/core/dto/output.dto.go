package dto

type Output struct {
	WattPerHourForLastMinute float64
	WattPerHourForThisHour   int64
	WattConsumedToday        int64
	EuroSpentToday           float64
	RatioLowTariffToday      float64
}

func NewOutput(dailySumLow int64, dailySumHigh int64, hourlySum int64, hourlyNbIndices int64, consumptionSinceLastTime int64, minutesSinceTheLastIndice float64) *Output {

	HIGH_TARIFF_PRICE := 0.0001963
	LOW_TARIFF_PRICE := 0.0001457

	ratioLowTariffToday := 1.0
	if dailySumHigh > 0 {
		ratioLowTariffToday = float64(dailySumLow) / float64(dailySumHigh+dailySumLow)
	}

	return &Output{
		WattPerHourForLastMinute: float64(consumptionSinceLastTime) * 60 / minutesSinceTheLastIndice,
		WattPerHourForThisHour:   hourlySum * 60 / hourlyNbIndices,
		WattConsumedToday:        dailySumHigh + dailySumLow,
		EuroSpentToday:           float64(dailySumHigh)*HIGH_TARIFF_PRICE + float64(dailySumLow)*LOW_TARIFF_PRICE,
		RatioLowTariffToday:      ratioLowTariffToday,
	}
}
