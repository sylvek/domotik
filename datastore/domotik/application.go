package main

import (
	"time"

	"github.com/sylvek/domotik/datastore/domain"
	"github.com/sylvek/domotik/datastore/port"
)

type Application struct {
	input                  domain.Input
	state                  domain.State
	stateRepository        port.StateRepository
	logRepository          port.LogRepository
	notificationRepository port.NotificationRepository
}

func (a *Application) AddLog(l domain.Log) error {
	if err := a.logRepository.Store(l); err != nil {
		return err
	}

	a.input.UpdateFromLog(l)

	return nil
}

func (a *Application) Process() error {
	if a.input.HasIndice() {
		t := time.Now()

		consumptionSinceLastTime, minutesSinceTheLastIndice := a.state.GetConsumptionSinceLastTime(t.Unix(), a.input.GetIndice())

		a.state.ProcessNewDay(t)
		a.state.ProcessNewHour(t)

		a.state.IncHourlyConsumption(consumptionSinceLastTime)
		a.state.IncDailyConsumption(consumptionSinceLastTime, a.input.IsLowTariff())

		output := *domain.NewOutput(
			a.state.DailySumLow,
			a.state.DailySumHigh,
			a.state.HourlySum,
			a.state.HourlyNbIndices,
			consumptionSinceLastTime,
			minutesSinceTheLastIndice)

		if err := a.stateRepository.Store(a.state); err != nil {
			return err
		}
		if err := a.notificationRepository.Notify(output); err != nil {
			return err
		}

		a.input.ResetIndice()
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
		input:                  *domain.NewInput(),
		state:                  state,
		stateRepository:        stateRepository,
		logRepository:          logRepository,
		notificationRepository: notificationRepository,
	}, nil
}
