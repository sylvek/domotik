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

func TestNotifyDayEnded(t *testing.T) {
	repository := &fakeRepository{values: map[string]float64{DAILY_POWER: 12345}}
	application, _ := NewApplication(&fakeState{}, repository, repository, NewSettings(nil))
	notifier := &fakeNotifier{}

	// nothing to notify while the day is running
	application.AddLog([]string{"sensors", "linky", "watt"}, 1)
	if err := application.NotifyDayEnded(notifier); err != nil || len(notifier.events) != 0 {
		t.Fatalf("unexpected notification %+v %v", notifier.events, err)
	}

	// the day ends: a failure is retried, a success is notified once
	repository.closeDay = true
	application.AddLog([]string{"sensors", "linky", "watt"}, 1)
	repository.closeDay = false

	notifier.err = errors.New("unreachable")
	if err := application.NotifyDayEnded(notifier); err == nil {
		t.Fatal("an error was expected")
	}
	notifier.err = nil
	application.NotifyDayEnded(notifier)
	application.NotifyDayEnded(notifier)
	if len(notifier.events) != 1 {
		t.Fatalf("expected one notification, got %+v", notifier.events)
	}
	event := notifier.events[0]
	if event.WattConsumed != 12345 || event.Timestamp == 0 || event.Timestamp%86400 != 0 {
		t.Errorf("unexpected event %+v", event)
	}

	// after too many failures, the notification is given up
	repository.closeDay = true
	application.AddLog([]string{"sensors", "linky", "watt"}, 1)
	repository.closeDay = false
	notifier.err = errors.New("unreachable")
	for i := 0; i < DAY_ENDED_MAX_ATTEMPTS; i++ {
		if err := application.NotifyDayEnded(notifier); err == nil {
			t.Fatal("an error was expected")
		}
	}
	if err := application.NotifyDayEnded(notifier); err != nil {
		t.Errorf("no more attempt was expected: %v", err)
	}
}
