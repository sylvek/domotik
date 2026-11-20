package core

import (
	"sort"
	"time"

	"github.com/sylvek/domotik/datastore/core/dto"
	"github.com/sylvek/domotik/datastore/core/model"
	"github.com/sylvek/domotik/datastore/core/port"
)

const (
	DAY                  = int64(24 * 60 * 60)
	TEMPERATURE_SLOT     = int64(5 * 60)
	DAILY_TEMP_OUTSIDE   = "daily_temp_outside"
	DAILY_TEMP_INSIDE    = "daily_temp_inside"
	DAILY_POWER          = "daily_power_consumption"
	DAILY_RATE           = "daily_rate_consumption"
	DAILY_COST           = "daily_cost_consumption"
	UNIT_TEMP            = "temp"
	UNIT_DEWPOINT        = "dewpoint"
	UNIT_WATT            = "watt"
	UNIT_STATE           = "state"
	UNIT_LITER           = "liter"
	QUARTER_DAYS         = 90
	HISTORY_YEARS        = 5
	LOW_TARIFF_STATE     = 0.0
	LINKY                = "linky"
	MEAN_PER_HOUR        = "meanPerHour"
	MEAN_PER_MINUTE      = "meanPerMinute"
	MEAN_PER_HOUR_LABEL  = "moyenne"
	MEAN_PER_MIN_LABEL   = "consommation"
	OUTSIDE_NOW_LABEL    = "now"
	OUTSIDE_DEW_LABEL    = "dew"
	OUTSIDE_SERIES_LABEL = "outside"
	INSIDE_SERIES_LABEL  = "inside"
)

func (a *Application) GetToday(presenter port.DashboardPresenter) error {
	now := time.Now()
	from := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	to := from.AddDate(0, 0, 1)

	summarize := a.summarize()
	today := dto.Today{
		From:                from.Unix(),
		To:                  to.Unix(),
		WattConsumedToday:   summarize.WattConsumedToday,
		EuroSpentToday:      summarize.EuroSpentToday,
		RatioLowTariffToday: summarize.RatioLowTariffToday,
		Inside:              []dto.Reading{},
		Outside:             []dto.Reading{},
		Temperatures:        []dto.Series{},
		Power:               []dto.Series{},
	}

	state, err := a.measureRepository.Last(LINKY, UNIT_STATE, today.From, today.To)
	if err != nil {
		return err
	}
	if state != nil {
		lowTariff := state.Value == LOW_TARIFF_STATE
		today.LowTariff = &lowTariff
	}

	watt, err := a.measureRepository.Last(LINKY, UNIT_WATT, today.From, today.To)
	if err != nil {
		return err
	}
	if watt != nil {
		today.WattNow = &watt.Value
	}

	temperatures, err := a.measureRepository.Series("", UNIT_TEMP, today.From, today.To)
	if err != nil {
		return err
	}
	for _, name := range a.sortedSensors(temperatures) {
		points := temperatures[name]
		last := dto.Reading{Label: a.settings.label(name), Value: points[len(points)-1].Value}
		if name == a.settings.OutsideSensor {
			last.Label = OUTSIDE_NOW_LABEL
			today.Outside = append(today.Outside, last)
		} else {
			today.Inside = append(today.Inside, last)
		}
		today.Temperatures = append(today.Temperatures,
			dto.NewSeries(a.settings.label(name), model.Downsample(points, TEMPERATURE_SLOT)))
	}

	dew, err := a.measureRepository.Last(a.settings.OutsideSensor, UNIT_DEWPOINT, today.From, today.To)
	if err != nil {
		return err
	}
	if dew != nil {
		today.Outside = append(today.Outside, dto.Reading{Label: OUTSIDE_DEW_LABEL, Value: dew.Value})
	}

	for _, s := range []struct{ name, label string }{
		{MEAN_PER_HOUR, MEAN_PER_HOUR_LABEL},
		{MEAN_PER_MINUTE, MEAN_PER_MIN_LABEL}} {
		series, err := a.measureRepository.Series(s.name, UNIT_WATT, today.From, today.To)
		if err != nil {
			return err
		}
		today.Power = append(today.Power, dto.NewSeries(s.label, series[s.name]))
	}

	return presenter.PresentToday(today)
}

func (a *Application) GetQuarter(presenter port.DashboardPresenter) error {
	now := time.Now()
	from, to := now.AddDate(0, 0, -QUARTER_DAYS).Unix(), now.Unix()

	daily := map[string][]model.Point{}
	for _, name := range []string{DAILY_TEMP_OUTSIDE, DAILY_TEMP_INSIDE, DAILY_POWER, DAILY_RATE, DAILY_COST} {
		points, err := a.measureRepository.Daily(name, from, to)
		if err != nil {
			return err
		}
		daily[name] = points
	}
	water, err := a.measureRepository.DailySum(UNIT_LITER, from, to)
	if err != nil {
		return err
	}

	rate := make([]model.Point, len(daily[DAILY_RATE]))
	for i, p := range daily[DAILY_RATE] {
		rate[i] = model.Point{TS: p.TS, Value: 100 * p.Value}
	}
	waterCost := model.GroupSum(water, model.ByWeek)
	for i := range waterCost {
		waterCost[i].Value *= a.settings.WaterEuroByLiter
	}

	return presenter.PresentQuarter(dto.Quarter{
		From: from,
		To:   to,
		Temperatures: []dto.Series{
			dto.NewSeries(OUTSIDE_SERIES_LABEL, daily[DAILY_TEMP_OUTSIDE]),
			dto.NewSeries(INSIDE_SERIES_LABEL, daily[DAILY_TEMP_INSIDE])},
		Consumption:      dto.NewSeries("consumption", daily[DAILY_POWER]),
		Water:            dto.NewSeries("water", water),
		LowTariffRate:    dto.NewSeries("low tariff rate", rate),
		WeeklyCost:       dto.NewBuckets(model.GroupSum(daily[DAILY_COST], model.ByWeek)),
		WeeklyWaterCost:  dto.NewBuckets(waterCost),
		WaterEuroByLiter: a.settings.WaterEuroByLiter,
	})
}

func (a *Application) GetHistory(presenter port.DashboardPresenter) error {
	now := time.Now()
	from, to := now.AddDate(-HISTORY_YEARS, 0, 0).Unix(), now.Unix()

	// yearly counters are computed on the whole history
	allConsumption, err := a.measureRepository.Daily(DAILY_POWER, 0, to)
	if err != nil {
		return err
	}
	allTemperature, err := a.measureRepository.Daily(DAILY_TEMP_OUTSIDE, 0, to)
	if err != nil {
		return err
	}

	consumptionByDay := map[int64]float64{}
	for _, p := range allConsumption {
		consumptionByDay[p.TS] = p.Value
	}
	highConsumptionOnCold := model.GroupCount(allTemperature, model.ByYear, func(p model.Point) bool {
		return p.Value < a.settings.ColdDayTemperature && consumptionByDay[p.TS] > a.settings.HighConsumptionWatt
	})
	hotDays := model.GroupCount(allTemperature, model.ByYear, func(p model.Point) bool {
		return p.Value > a.settings.HotDayTemperature
	})

	consumption := since(allConsumption, from)

	return presenter.PresentHistory(dto.History{
		From:                  from,
		To:                    to,
		Consumption:           dto.NewSeries("consumption", consumption),
		Temperature:           dto.NewSeries(OUTSIDE_SERIES_LABEL, since(allTemperature, from)),
		MonthlyConsumption:    dto.NewBuckets(model.GroupSum(consumption, model.ByMonth)),
		YearlyConsumption:     dto.NewBuckets(model.GroupSum(consumption, model.ByYear)),
		HighConsumptionOnCold: dto.NewBuckets(highConsumptionOnCold),
		HotDays:               dto.NewBuckets(hotDays),
		ColdDayTemperature:    a.settings.ColdDayTemperature,
		HighConsumptionWatt:   a.settings.HighConsumptionWatt,
		HotDayTemperature:     a.settings.HotDayTemperature,
	})
}

// sortedSensors returns the outside sensor first, then the others sorted by label.
func (a *Application) sortedSensors(series map[string][]model.Point) []string {
	names := []string{}
	for name, points := range series {
		if len(points) > 0 {
			names = append(names, name)
		}
	}
	sort.Slice(names, func(i, j int) bool {
		if (names[i] == a.settings.OutsideSensor) != (names[j] == a.settings.OutsideSensor) {
			return names[i] == a.settings.OutsideSensor
		}
		return a.settings.label(names[i]) < a.settings.label(names[j])
	})
	return names
}

func since(points []model.Point, from int64) []model.Point {
	i := sort.Search(len(points), func(i int) bool { return points[i].TS >= from })
	return points[i:]
}
