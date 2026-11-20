package core

import (
	"errors"
	"testing"

	"github.com/sylvek/domotik/datastore/core/dto"
	"github.com/sylvek/domotik/datastore/core/model"
)

type fakeRepository struct {
	closeDay bool
	values   map[string]float64
}

func (f *fakeRepository) Close()                {}
func (f *fakeRepository) Store(model.Log) error { return nil }
func (f *fakeRepository) CloseDay() bool        { return f.closeDay }
func (f *fakeRepository) Retrieve() (model.State, error) {
	return *model.NewState(), nil
}
func (f *fakeRepository) Last(string, string, int64, int64) (*model.Point, error) {
	return nil, nil
}
func (f *fakeRepository) Series(string, string, int64, int64) (map[string][]model.Point, error) {
	return nil, nil
}
func (f *fakeRepository) Daily(name string, from int64, to int64) ([]model.Point, error) {
	if value, ok := f.values[name]; ok {
		return []model.Point{{TS: from, Value: value}}, nil
	}
	return nil, nil
}
func (f *fakeRepository) DailySum(string, int64, int64) ([]model.Point, error) { return nil, nil }

type fakeState struct{ fakeRepository }

func (f *fakeState) Store(model.State) error { return nil }

type fakeNotifier struct {
	events []dto.DayEnded
	err    error
}

func (f *fakeNotifier) Notify(event dto.DayEnded) error {
	if f.err != nil {
		return f.err
	}
	f.events = append(f.events, event)
	return nil
}

type fakePresenter struct{}

func (f *fakePresenter) Present(dto.Output) error { return nil }

func TestProcessNotifiesDayEnded(t *testing.T) {
	repository := &fakeRepository{values: map[string]float64{DAILY_POWER: 12345}}
	application, _ := NewApplication(&fakeState{}, repository, repository, NewSettings(nil))
	notifier := &fakeNotifier{}
	presenter := &fakePresenter{}

	// nothing to notify while the day is running
	application.AddLog([]string{"sensors", "linky", "watt"}, 1)
	if err := application.Process(presenter, notifier); err != nil || len(notifier.events) != 0 {
		t.Fatalf("unexpected notification %+v %v", notifier.events, err)
	}

	// the day ends: it is notified once
	repository.closeDay = true
	application.AddLog([]string{"sensors", "linky", "watt"}, 1)
	repository.closeDay = false
	application.Process(presenter, notifier)
	application.Process(presenter, notifier)
	if len(notifier.events) != 1 {
		t.Fatalf("expected one notification, got %+v", notifier.events)
	}
	event := notifier.events[0]
	if event.WattConsumed != 12345 || event.Timestamp == 0 || event.Timestamp%86400 != 0 {
		t.Errorf("unexpected event %+v", event)
	}

	// a day without any consumption is not notified
	delete(repository.values, DAILY_POWER)
	repository.closeDay = true
	application.AddLog([]string{"sensors", "linky", "watt"}, 1)
	repository.closeDay = false
	if err := application.Process(presenter, notifier); err != nil || len(notifier.events) != 1 {
		t.Fatalf("unexpected notification %+v %v", notifier.events, err)
	}

	// a notifier failure is reported
	repository.values[DAILY_POWER] = 1
	repository.closeDay = true
	application.AddLog([]string{"sensors", "linky", "watt"}, 1)
	notifier.err = errors.New("failure")
	if err := application.Process(presenter, notifier); err == nil {
		t.Error("an error was expected")
	}
}
