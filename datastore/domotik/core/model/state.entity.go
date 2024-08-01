package model

type State struct {
	LastIndice                int64
	LastIndiceTS              int64
	CurrentDay                int
	CurrentHour               int
	DailySumHigh              int64
	DailySumLow               int64
	HourlySum                 int64
	HourlyNbIndices           int64
	consumptionSinceLastTime  int64
	minutesSinceTheLastIndice float64
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

func (s *State) GetConsumptionSinceLastTime() int64 {
	return s.consumptionSinceLastTime
}

func (s *State) GetMinutesSinceTheLastIndice() float64 {
	return s.minutesSinceTheLastIndice
}

func (s *State) setConsumptionSinceLastTime(now int64, newIndice int64) {
	if s.LastIndice > 0 && s.LastIndiceTS > 0 {
		diff := newIndice - s.LastIndice
		if diff > 0 {
			s.consumptionSinceLastTime = diff
			s.minutesSinceTheLastIndice = float64(now-s.LastIndiceTS) / 60
		}
	}
	s.LastIndice = newIndice
	s.LastIndiceTS = now
}

func (s *State) ProcessNewDay(currentDay int) {
	if currentDay != s.CurrentDay {
		s.CurrentDay = currentDay
		s.DailySumHigh = 0
		s.DailySumLow = 0
	}
}

func (s *State) ProcessNewHour(currentHour int) {
	if currentHour != s.CurrentHour {
		s.CurrentHour = currentHour
		s.HourlyNbIndices = 0
		s.HourlySum = 0
	}
}

func (s *State) incHourlyConsumption(wattConsumedDuringBuffering int64) {
	s.HourlySum += wattConsumedDuringBuffering
	s.HourlyNbIndices += 1
}

func (s *State) incDailyConsumption(wattConsumedDuringBuffering int64, lowTariff bool) {
	if lowTariff {
		s.DailySumLow += wattConsumedDuringBuffering
	} else {
		s.DailySumHigh += wattConsumedDuringBuffering
	}
}

func (s *State) ProcessLastIndice(now int64, measure Measure) {
	newIndice := measure.GetLastIndex()
	isLowTariff := measure.IsLowTariff()
	s.setConsumptionSinceLastTime(now, newIndice)
	s.incHourlyConsumption(s.consumptionSinceLastTime)
	s.incDailyConsumption(s.consumptionSinceLastTime, isLowTariff)
}
