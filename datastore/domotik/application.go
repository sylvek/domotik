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
		newState, output := domain.GenerateStatistics(time.Now(), a.state, a.input)

		if err := a.stateRepository.Store(newState); err != nil {
			return err
		}
		if err := a.notificationRepository.Notify(output); err != nil {
			return err
		}

		a.input.ResetIndice()
		a.state = newState
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
