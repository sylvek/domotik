package dto

import "github.com/sylvek/domotik/datastore/core/model"

type Output struct {
	WattPerHourForLastMinute float64
	WattPerHourForThisHour   int64
	WattConsumedToday        int64
	EuroSpentToday           float64
	RatioLowTariffToday      float64
}

func NewOutputFromState(state model.State) *Output {

	HIGH_TARIFF_PRICE := 0.0001963
	LOW_TARIFF_PRICE := 0.0001457

	dailySumHigh := state.DailySumHigh
	dailySumLow := state.DailySumLow
	hourlySum := state.HourlySum
	hourlyNbIndices := state.HourlyNbIndices
	consumptionSinceLastTime := state.GetConsumptionSinceLastTime()
	minutesSinceTheLastIndice := state.GetMinutesSinceTheLastIndice()

	wattPerHourForLastMinute := 0.0
	if minutesSinceTheLastIndice > 0 {
		wattPerHourForLastMinute = float64(consumptionSinceLastTime) * 60 / minutesSinceTheLastIndice
	}

	wattPerHourForThisHour := int64(0)
	if hourlyNbIndices > 0 {
		wattPerHourForThisHour = hourlySum * 60 / hourlyNbIndices
	}

	ratioLowTariffToday := 1.0
	if dailySumHigh > 0 {
		ratioLowTariffToday = float64(dailySumLow) / float64(dailySumHigh+dailySumLow)
	}

	return &Output{
		WattPerHourForLastMinute: wattPerHourForLastMinute,
		WattPerHourForThisHour:   wattPerHourForThisHour,
		WattConsumedToday:        dailySumHigh + dailySumLow,
		EuroSpentToday:           float64(dailySumHigh)*HIGH_TARIFF_PRICE + float64(dailySumLow)*LOW_TARIFF_PRICE,
		RatioLowTariffToday:      ratioLowTariffToday,
	}
}
