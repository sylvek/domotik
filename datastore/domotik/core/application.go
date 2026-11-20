package core

import (
	"sync"
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

	measureRepository port.MeasureRepository
	settings          Settings
	mutex             sync.RWMutex

	// day (midnight UTC) that has just ended and is not notified yet, 0 if none.
	endedDay         int64
	endedDayAttempts int
}

func (a *Application) AddLog(elements []string, value float64) error {
	log := *model.NewLog(elements, value)
	dayEnded := a.logRepository.CloseDay()
	if err := a.logRepository.Store(log); err != nil {
		return err
	}

	a.mutex.Lock()
	defer a.mutex.Unlock()
	if dayEnded {
		a.endedDay = time.Now().UTC().Truncate(24 * time.Hour).Unix()
		a.endedDayAttempts = 0
	}
	a.measure.UpdateFromLog(log, time.Now().Unix())

	return nil
}

// ConvertTasmota publishes the measures of a tasmota device as regular sensors logs.
func (a *Application) ConvertTasmota(telemetry model.TasmotaTelemetry, publisher port.LogPublisher) error {
	for _, log := range telemetry.Logs() {
		if err := publisher.Publish(log); err != nil {
			return err
		}
	}

	return nil
}

func (a *Application) Process(presenter port.SummarizePresenter) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.measure.HasNewIndex() {

		t := time.Now()

		a.state.ProcessNewDay(t.Day())
		a.state.ProcessNewHour(t.Hour())
		a.state.ProcessLastIndice(t.Unix(), a.measure)

		if err := a.stateRepository.Store(a.state); err != nil {
			return err
		}

		output := *dto.NewOutputFromState(a.state)

		if err := presenter.Present(output); err != nil {
			return err
		}

		a.measure.ResetIndex()
	}

	return nil
}

func (a *Application) GetSummarize(presenter port.SummarizePresenter) error {
	return presenter.Present(a.summarize())
}

func (a *Application) summarize() dto.Output {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	return *dto.NewOutputFromState(a.state)
}

func NewApplication(
	stateRepository port.StateRepository,
	logRepository port.LogRepository,
	measureRepository port.MeasureRepository,
	settings Settings) (*Application, error) {

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

		measureRepository: measureRepository,
		settings:          settings,
	}, nil
}
