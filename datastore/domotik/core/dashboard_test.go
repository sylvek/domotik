package core

import (
	"testing"

	"github.com/sylvek/domotik/datastore/core/dto"
	"github.com/sylvek/domotik/datastore/core/model"
)

type todayPresenter struct{ today dto.Today }

func (p *todayPresenter) PresentToday(today dto.Today) error { p.today = today; return nil }
func (p *todayPresenter) PresentQuarter(dto.Quarter) error   { return nil }
func (p *todayPresenter) PresentHistory(dto.History) error   { return nil }

func TestGetTodayFallsBackToTheCurrentDay(t *testing.T) {
	state := *model.NewState()
	state.DailySumHigh = 12000
	repository := &fakeRepository{state: state}
	application, _ := NewApplication(fakeStateRepository{repository}, fakeLogRepository{repository}, repository, NewSettings(nil))

	for daysAgo, expected := range map[int]int{0: 0, 1: 1, 2: 0, -1: 0} {
		presenter := &todayPresenter{}
		if err := application.GetToday(daysAgo, presenter); err != nil {
			t.Fatal(err)
		}
		if presenter.today.DaysAgo != expected || presenter.today.Live != (expected == 0) {
			t.Errorf("%d days ago: unexpected day %d (live: %v)", daysAgo, presenter.today.DaysAgo, presenter.today.Live)
		}
		if presenter.today.To-presenter.today.From < 23*3600 {
			t.Errorf("unexpected period %d - %d", presenter.today.From, presenter.today.To)
		}
		// the state only describes the current day
		if expected == 1 && presenter.today.WattConsumedToday != 0 {
			t.Errorf("the state must not be used for the day before, got %d", presenter.today.WattConsumedToday)
		}
	}
}
