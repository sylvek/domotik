package core

import (
	"testing"
	"time"

	"github.com/sylvek/domotik/datastore/core/dto"
	"github.com/sylvek/domotik/datastore/core/model"
)

type fakeRepository struct {
	state model.State
}

func (f *fakeRepository) Close()                                            {}
func (f *fakeRepository) Retrieve() (model.State, error)                    { return f.state, nil }
func (f *fakeRepository) Daily(string, int64, int64) ([]model.Point, error) { return nil, nil }
func (f *fakeRepository) DailyWater(int64, int64) ([]model.Point, error) {
	return nil, nil
}
func (f *fakeRepository) Last(string, string, int64, int64) (*model.Point, error) {
	return nil, nil
}
func (f *fakeRepository) Series(string, string, int64, int64) (map[string][]model.Point, error) {
	return nil, nil
}

type fakeStateRepository struct{ *fakeRepository }

func (f fakeStateRepository) Store(model.State) error { return nil }

type fakeLogRepository struct{ *fakeRepository }

func (f fakeLogRepository) Store(model.Log) error { return nil }

type fakePresenter struct {
	days []dto.DayEnded
}

func (f *fakePresenter) Present(dto.Output) error { return nil }

type fakeDayEndedPresenter struct{ *fakePresenter }

func (f fakeDayEndedPresenter) Present(day dto.DayEnded) error {
	f.days = append(f.days, day)
	return nil
}

func TestProcessPresentsTheEndedDay(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	// yesterday, 12 kWh were consumed
	state := *model.NewState()
	state.CurrentDay = yesterday.Day()
	state.CurrentHour = yesterday.Hour()
	state.LastIndice = 100000
	state.LastIndiceTS = yesterday.Unix()
	state.DailySumHigh = 2700
	state.DailySumLow = 9300

	repository := &fakeRepository{state: state}
	application, _ := NewApplication(fakeStateRepository{repository}, fakeLogRepository{repository}, repository, NewSettings(nil))
	presenter := &fakePresenter{}

	// without any new index, nothing happens
	application.Process(presenter, fakeDayEndedPresenter{presenter})
	if len(presenter.days) != 0 {
		t.Fatalf("unexpected day %+v", presenter.days)
	}

	// the first index of the new day ends the previous one
	application.AddLog([]string{"sensors", "linky", "indice"}, 100500)
	application.Process(presenter, fakeDayEndedPresenter{presenter})
	if len(presenter.days) != 1 {
		t.Fatalf("expected one day, got %+v", presenter.days)
	}
	expected := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, time.UTC).Unix()
	if presenter.days[0].WattConsumed != 12000 || presenter.days[0].Timestamp != expected {
		t.Errorf("unexpected day %+v", presenter.days[0])
	}
	if summarize := application.summarize(); summarize.WattConsumedToday != 500 {
		t.Errorf("the new day should start with the last consumption, got %d", summarize.WattConsumedToday)
	}

	// the day is presented once
	application.AddLog([]string{"sensors", "linky", "indice"}, 100600)
	application.Process(presenter, fakeDayEndedPresenter{presenter})
	if len(presenter.days) != 1 {
		t.Errorf("expected one day, got %+v", presenter.days)
	}
}

func TestProcessDoesNotPresentAnEmptyDay(t *testing.T) {
	repository := &fakeRepository{state: *model.NewState()}
	application, _ := NewApplication(fakeStateRepository{repository}, fakeLogRepository{repository}, repository, NewSettings(nil))
	presenter := &fakePresenter{}

	application.AddLog([]string{"sensors", "linky", "indice"}, 100500)
	application.Process(presenter, fakeDayEndedPresenter{presenter})
	if len(presenter.days) != 0 {
		t.Errorf("unexpected day %+v", presenter.days)
	}
}
