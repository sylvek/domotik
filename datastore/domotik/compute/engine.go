package compute

import (
	"encoding/json"
	"os"
	"time"
)

type input struct {
	lowTariff bool
	indice    int64
}

type state struct {
	LastIndice      int64
	LastIndiceTS    int64
	CurrentDay      int
	CurrentHour     int
	DailySumHigh    int64
	DailySumLow     int64
	HourlySum       int64
	HourlyNbIndices int64
}

type RuleEngineClient struct {
	path   string
	ticker *time.Ticker
	input  *input
	state  *state
	output chan Output
	done   chan bool
}

func NewRuleEngineClient(path string) Rule {
	state := &state{
		LastIndice:      0,
		CurrentDay:      0,
		CurrentHour:     0,
		DailySumHigh:    0,
		DailySumLow:     0,
		HourlySum:       0,
		HourlyNbIndices: 0,
	}
	data, err := os.ReadFile(path + "/state.json")
	if err == nil {
		json.Unmarshal(data, &state)
	}

	instance := &RuleEngineClient{
		path:   path,
		input:  &input{lowTariff: false, indice: 0},
		ticker: time.NewTicker(time.Minute),
		state:  state,
		output: make(chan Output),
		done:   make(chan bool)}
	go instance.start()
	return instance
}

func (g *RuleEngineClient) start() {
	for {
		select {
		case <-g.ticker.C:
			if g.input.indice > 0 {
				g.output <- g.updateStateAndGenerateOutput()
				g.input.indice = 0
				g.synchronization()
			}
		case <-g.done:
			return
		}
	}
}

func (g *RuleEngineClient) synchronization() {
	s, _ := json.Marshal(g.state)
	os.WriteFile(g.path+"/state.json", s, 0644)
}

func (g *RuleEngineClient) updateStateAndGenerateOutput() Output {
	t := time.Now()
	now := t.Unix()

	lastIndice := g.state.LastIndice
	lastEpoch := g.state.LastIndiceTS
	newIndice := g.input.indice

	// calculate watt consumed between 2 Ticks
	wattConsumedDuringBuffering := int64(0)
	if lastIndice > 0 {
		wattConsumedDuringBuffering = newIndice - lastIndice
	}
	g.state.LastIndice = newIndice

	minutesSinceTheLastIndice := float64(1)
	if wattConsumedDuringBuffering > 0 {
		if lastEpoch > 0 {
			minutesSinceTheLastIndice = float64(now-lastEpoch) / 60
		}
		g.state.LastIndiceTS = now
	}

	// if new day -> clear daily state
	if t.Day() != g.state.CurrentDay {
		g.state.CurrentDay = t.Day()
		g.state.DailySumHigh = 0
		g.state.DailySumLow = 0
	}

	// if new hour -> clear hourly state
	if t.Hour() != g.state.CurrentHour {
		g.state.CurrentHour = t.Hour()
		g.state.HourlyNbIndices = 0
		g.state.HourlySum = 0
	}

	// calculate mean per hour
	g.state.HourlySum += wattConsumedDuringBuffering
	g.state.HourlyNbIndices += 1

	// calculate ratio low/high and €
	if g.input.lowTariff {
		g.state.DailySumLow += wattConsumedDuringBuffering
	} else {
		g.state.DailySumHigh += wattConsumedDuringBuffering
	}

	ratioLowTariffToday := 1.0
	if g.state.DailySumHigh > 0 {
		ratioLowTariffToday = float64(g.state.DailySumLow) / float64(g.state.DailySumHigh+g.state.DailySumLow)
	}

	return Output{
		WattPerHourForLastMinute: float64(wattConsumedDuringBuffering) * 60 / minutesSinceTheLastIndice,
		WattPerHourForThisHour:   g.state.HourlySum * 60 / g.state.HourlyNbIndices,
		WattConsumedToday:        g.state.DailySumHigh + g.state.DailySumLow,
		EuroSpentToday:           float64(g.state.DailySumHigh)*0.0001963 + float64(g.state.DailySumLow)*0.0001457,
		RatioLowTariffToday:      ratioLowTariffToday,
	}
}

func (g *RuleEngineClient) Tick() chan Output {
	return g.output
}

func (g *RuleEngineClient) Stop() {
	g.ticker.Stop()
	g.synchronization()
	g.done <- true
}

func (g *RuleEngineClient) SetLowTariffState(state bool) {
	g.input.lowTariff = state
}

func (g *RuleEngineClient) SetIndice(indice int64) {
	g.input.indice = indice
}
