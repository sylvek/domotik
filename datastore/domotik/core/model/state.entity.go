package model

import "time"

type State struct {
	LastIndice      int64
	LastIndiceTS    int64
	CurrentDay      int
	CurrentHour     int
	DailySumHigh    int64
	DailySumLow     int64
	HourlySum       int64
	HourlyNbIndices int64
}

func NewState() *State {
	return &State{
		LastIndice:      0,
		CurrentDay:      0,
		CurrentHour:     0,
		DailySumHigh:    0,
		DailySumLow:     0,
		HourlySum:       0,
		HourlyNbIndices: 0,
	}
}

func (s *State) GetConsumptionSinceLastTime(now int64, newIndice int64) (int64, float64) {
	consumptionSinceLastTime := int64(0)
	if s.LastIndice > 0 {
		consumptionSinceLastTime = newIndice - s.LastIndice
	}
	s.LastIndice = newIndice

	minutesSinceTheLastIndice := float64(1)
	if consumptionSinceLastTime > 0 {
		if s.LastIndiceTS > 0 {
			minutesSinceTheLastIndice = float64(now-s.LastIndiceTS) / 60
		}
	}

	return consumptionSinceLastTime, minutesSinceTheLastIndice
}

func (s *State) ProcessNewDay(t time.Time) {
	if t.Day() != s.CurrentDay {
		s.CurrentDay = t.Day()
		s.DailySumHigh = 0
		s.DailySumLow = 0
	}
}

func (s *State) ProcessNewHour(t time.Time) {
	if t.Hour() != s.CurrentHour {
		s.CurrentHour = t.Hour()
		s.HourlyNbIndices = 0
		s.HourlySum = 0
	}
}

func (s *State) IncHourlyConsumption(wattConsumedDuringBuffering int64) {
	s.HourlySum += wattConsumedDuringBuffering
	s.HourlyNbIndices += 1
}

func (s *State) IncDailyConsumption(wattConsumedDuringBuffering int64, lowTariff bool) {
	if lowTariff {
		s.DailySumLow += wattConsumedDuringBuffering
	} else {
		s.DailySumHigh += wattConsumedDuringBuffering
	}
}
