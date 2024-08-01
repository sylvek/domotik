package core

import (
	"time"

	"github.com/sylvek/domotik/datastore/core/dto"
	"github.com/sylvek/domotik/datastore/core/model"
	"github.com/sylvek/domotik/datastore/core/port"
)

type Application struct {
	measure         model.Measure
	state           model.State
	stateRepository port.StateRepository
	logRepository   port.LogRepository
}

func (a *Application) AddLog(elements []string, value float64) error {
	log := *model.NewLog(elements, value)
	if err := a.logRepository.Store(log); err != nil {
		return err
	}

	a.measure.UpdateFromLog(log)

	return nil
}

func (a *Application) Process(presenter port.SummarizePresenter) error {
	if a.measure.HasNewIndex() {

		t := time.Now()
		a.state.ProcessNewDay(t)
		a.state.ProcessNewHour(t)

		output := a.summarizeSince(t.Unix())
		if err := presenter.Present(output); err != nil {
			return err
		}

		if err := a.stateRepository.Store(a.state); err != nil {
			return err
		}

		a.measure.ResetIndex()
	}

	return nil
}

func (a *Application) GetSummarize(presenter port.SummarizePresenter) error {
	output := a.summarizeSince(a.state.LastIndiceTS)
	if err := presenter.Present(output); err != nil {
		return err
	}

	return nil
}

func NewApplication(
	stateRepository port.StateRepository,
	logRepository port.LogRepository) (*Application, error) {

	measure := *model.NewMeasure()

	state, err := stateRepository.Retrieve()
	if err != nil {
		return nil, err
	}

	return &Application{
		measure:         measure,
		state:           state,
		stateRepository: stateRepository,
		logRepository:   logRepository,
	}, nil
}

func (a *Application) summarizeSince(timestamp int64) dto.Output {

	doChanges := a.state.LastIndiceTS < timestamp

	consumptionSinceLastTime, minutesSinceTheLastIndice := a.state.GetConsumptionSinceLastTime(timestamp, a.measure.GetLastIndex())

	if doChanges {
		a.state.LastIndiceTS = timestamp
		a.state.IncHourlyConsumption(consumptionSinceLastTime)
		a.state.IncDailyConsumption(consumptionSinceLastTime, a.measure.IsLowTariff())
	}

	return *dto.NewOutput(
		a.state.DailySumLow,
		a.state.DailySumHigh,
		a.state.HourlySum,
		a.state.HourlyNbIndices,
		consumptionSinceLastTime,
		minutesSinceTheLastIndice)
}
