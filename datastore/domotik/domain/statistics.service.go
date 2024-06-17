package domain

import (
	"time"
)

func GenerateStatistics(t time.Time, state State, input Input) (State, Output) {

	newState := CopyState(state)

	consumptionSinceLastTime, minutesSinceTheLastIndice := newState.GetConsumptionSinceLastTime(t.Unix(), input.indice)

	newState.ProcessNewDay(t)
	newState.ProcessNewHour(t)

	newState.IncHourlyConsumption(consumptionSinceLastTime)
	newState.IncDailyConsumption(consumptionSinceLastTime, input.lowTariff)

	ratioLowTariffToday := newState.GetLowTariffRate()

	return *newState, *NewOutput(*newState, consumptionSinceLastTime, minutesSinceTheLastIndice, ratioLowTariffToday)
}
