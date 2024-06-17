package domain

import (
	"time"

	"github.com/sylvek/domotik/datastore/domain/model"
	"github.com/sylvek/domotik/datastore/port"
)

type Application struct {
	input                  model.Input
	state                  model.State
	stateRepository        port.StateRepository
	logRepository          port.LogRepository
	notificationRepository port.NotificationRepository
}

func (a *Application) AddLog(l model.Log) error {
	if err := a.logRepository.Store(l); err != nil {
		return err
	}

	if l.Name == "linky" && l.Unit == "indice" {
		a.input.Indice = int64(l.Value)
	}

	if l.Name == "linky" && l.Unit == "state" {
		a.input.LowTariff = l.Value == 0.0
	}

	return nil
}

func (a *Application) Process() error {
	if a.input.Indice > 0 {
		newState, output := execute(time.Now(), a.state, a.input)
		a.input.Indice = 0
		a.state = newState

		if err := a.stateRepository.Store(newState); err != nil {
			return err
		}
		if err := a.notificationRepository.Notify(output); err != nil {
			return err
		}
	}

	return nil
}

func NewApplication(
	stateRepository port.StateRepository,
	logRepository port.LogRepository,
	notificationRepository port.NotificationRepository) (*Application, error) {

	state, err := stateRepository.Retrieve()
	if err != nil {
		return nil, err
	}

	return &Application{
		input:                  model.Input{LowTariff: false, Indice: 0},
		state:                  state,
		stateRepository:        stateRepository,
		logRepository:          logRepository,
		notificationRepository: notificationRepository,
	}, nil
}

func execute(t time.Time, state model.State, input model.Input) (model.State, model.Output) {
	now := t.Unix()

	lastIndice := state.LastIndice
	lastEpoch := state.LastIndiceTS
	newIndice := input.Indice

	// calculate watt consumed between 2 Ticks
	wattConsumedDuringBuffering := int64(0)
	if lastIndice > 0 {
		wattConsumedDuringBuffering = newIndice - lastIndice
	}
	state.LastIndice = newIndice

	minutesSinceTheLastIndice := float64(1)
	if wattConsumedDuringBuffering > 0 {
		if lastEpoch > 0 {
			minutesSinceTheLastIndice = float64(now-lastEpoch) / 60
		}
		state.LastIndiceTS = now
	}

	// if new day -> clear daily state
	if t.Day() != state.CurrentDay {
		state.CurrentDay = t.Day()
		state.DailySumHigh = 0
		state.DailySumLow = 0
	}

	// if new hour -> clear hourly state
	if t.Hour() != state.CurrentHour {
		state.CurrentHour = t.Hour()
		state.HourlyNbIndices = 0
		state.HourlySum = 0
	}

	// calculate mean per hour
	state.HourlySum += wattConsumedDuringBuffering
	state.HourlyNbIndices += 1

	// calculate ratio low/high and €
	if input.LowTariff {
		state.DailySumLow += wattConsumedDuringBuffering
	} else {
		state.DailySumHigh += wattConsumedDuringBuffering
	}

	ratioLowTariffToday := 1.0
	if state.DailySumHigh > 0 {
		ratioLowTariffToday = float64(state.DailySumLow) / float64(state.DailySumHigh+state.DailySumLow)
	}

	return state, model.Output{
		WattPerHourForLastMinute: float64(wattConsumedDuringBuffering) * 60 / minutesSinceTheLastIndice,
		WattPerHourForThisHour:   state.HourlySum * 60 / state.HourlyNbIndices,
		WattConsumedToday:        state.DailySumHigh + state.DailySumLow,
		EuroSpentToday:           float64(state.DailySumHigh)*0.0001963 + float64(state.DailySumLow)*0.0001457,
		RatioLowTariffToday:      ratioLowTariffToday,
	}
}
