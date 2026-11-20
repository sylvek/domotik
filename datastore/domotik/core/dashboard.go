package core

import (
	"slices"
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
	HISTORY_YEARS        = 10
	LOW_TARIFF_STATE     = 0.0
	SUM_PER_DAY          = "sumPerDay"
	UNIT_EURO            = "euro"
	UNIT_RATIO           = "rate"
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

// QUARTER_DAYS are the periods the quarter dashboard can display, the first one is the default.
var QUARTER_DAYS = []int{90, 180}

// TODAY_DAYS_AGO are the days the today dashboard can display: the current one (0) and the day before (1).
// Raw measures are kept two days.
var TODAY_DAYS_AGO = []int{0, 1}

// GetToday presents a whole day, from midnight to midnight: the current one, or the day before.
// An unexpected number of days falls back to the current day.
func (a *Application) GetToday(daysAgo int, presenter port.DashboardPresenter) error {
	if !slices.Contains(TODAY_DAYS_AGO, daysAgo) {
		daysAgo = TODAY_DAYS_AGO[0]
	}
	live := daysAgo == 0

	now := time.Now()
	from := time.Date(now.Year(), now.Month(), now.Day()-daysAgo, 0, 0, 0, 0, now.Location())
	to := from.AddDate(0, 0, 1)

	// the current day is summarized by the state, a past one by its last published summary
	summarize := a.summarize()
	if !live {
		past, err := a.summarizePast(from.Unix(), to.Unix())
		if err != nil {
			return err
		}
		summarize = past
	}
	typicalDay := a.settings.typicalDay(summarize.WattConsumedToday)
	today := dto.Today{
		From:                from.Unix(),
		To:                  to.Unix(),
		DaysAgo:             daysAgo,
		AvailableDaysAgo:    TODAY_DAYS_AGO,
		Live:                live,
		WattConsumedToday:   summarize.WattConsumedToday,
		EuroSpentToday:      summarize.EuroSpentToday,
		RatioLowTariffToday: summarize.RatioLowTariffToday,
		TypicalDay:          typicalDay.Label,
		TypicalDayLevel:     typicalDay.Level,
		TypicalDayMaxWatt:   a.settings.TypicalDays[len(a.settings.TypicalDays)-1].Watt,
		Inside:              []dto.Reading{},
		Outside:             []dto.Reading{},
		Temperatures:        []dto.Series{},
		Power:               []dto.Series{},
	}

	// instant values only make sense for the current day
	if live {
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
	}

	temperatures, err := a.measureRepository.Series("", UNIT_TEMP, today.From, today.To)
	if err != nil {
		return err
	}
	for _, name := range a.sortedSensors(temperatures) {
		points := temperatures[name]
		last := dto.Reading{Label: a.settings.label(name), Value: points[len(points)-1].Value}
		// no instant reading for a past day
		if live && name == a.settings.OutsideSensor {
			last.Label = OUTSIDE_NOW_LABEL
			today.Outside = append(today.Outside, last)
		} else if live {
			today.Inside = append(today.Inside, last)
		}
		today.Temperatures = append(today.Temperatures,
			dto.NewSeries(a.settings.label(name), model.Downsample(points, TEMPERATURE_SLOT)))
	}

	dew, err := a.measureRepository.Last(a.settings.OutsideSensor, UNIT_DEWPOINT, today.From, today.To)
	if err != nil {
		return err
	}
	if dew != nil && live {
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

// summarizePast rebuilds the summary of a past day from the last values published during that day.
func (a *Application) summarizePast(from int64, to int64) (dto.Output, error) {
	output := dto.Output{}
	for unit, set := range map[string]func(float64){
		UNIT_WATT:  func(v float64) { output.WattConsumedToday = int64(v) },
		UNIT_EURO:  func(v float64) { output.EuroSpentToday = v },
		UNIT_RATIO: func(v float64) { output.RatioLowTariffToday = v },
	} {
		last, err := a.measureRepository.Last(SUM_PER_DAY, unit, from, to)
		if err != nil {
			return output, err
		}
		if last != nil {
			set(last.Value)
		}
	}
	return output, nil
}

// GetQuarter presents the last `days` days. An unexpected number of days falls back to the default one.
func (a *Application) GetQuarter(days int, presenter port.DashboardPresenter) error {
	if !slices.Contains(QUARTER_DAYS, days) {
		days = QUARTER_DAYS[0]
	}

	now := time.Now()
	from, to := now.AddDate(0, 0, -days), now
	// the same period, a year before
	previousFrom, previousTo := from.AddDate(-1, 0, 0).Unix(), to.AddDate(-1, 0, 0).Unix()

	daily := map[string][]model.Point{}
	previous := map[string][]model.Point{}
	for _, name := range []string{DAILY_TEMP_OUTSIDE, DAILY_TEMP_INSIDE, DAILY_POWER, DAILY_RATE, DAILY_COST} {
		points, err := a.measureRepository.Daily(name, previousFrom, to.Unix())
		if err != nil {
			return err
		}
		daily[name] = model.Between(points, from.Unix(), to.Unix())
		previous[name] = model.Between(points, previousFrom, previousTo)
	}
	allWater, err := a.measureRepository.DailyWater(previousFrom, to.Unix())
	if err != nil {
		return err
	}
	water := model.Between(allWater, from.Unix(), to.Unix())

	rate := make([]model.Point, len(daily[DAILY_RATE]))
	for i, p := range daily[DAILY_RATE] {
		rate[i] = model.Point{TS: p.TS, Value: 100 * p.Value}
	}
	waterCost := model.GroupSum(water, model.ByWeek)
	for i := range waterCost {
		waterCost[i].Value *= a.settings.WaterEuroByLiter
	}

	return presenter.PresentQuarter(dto.Quarter{
		From:             from.Unix(),
		To:               to.Unix(),
		Days:             days,
		AvailableDays:    QUARTER_DAYS,
		ConsumptionTrend: dto.NewTrend(model.NewTrend(daily[DAILY_POWER], previous[DAILY_POWER], model.Sum)),
		CostTrend:        dto.NewTrend(model.NewTrend(daily[DAILY_COST], previous[DAILY_COST], model.Sum)),
		WaterTrend:       dto.NewTrend(model.NewTrend(water, model.Between(allWater, previousFrom, previousTo), model.Sum)),
		TemperatureTrend: dto.NewTrend(model.NewTrend(daily[DAILY_TEMP_OUTSIDE], previous[DAILY_TEMP_OUTSIDE], model.Mean)),
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

	// yearly counters and trends are computed on the whole history
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

	// since the 1st of january, compared to last year at the same date
	newYear := time.Date(now.UTC().Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	yearToDate := model.NewTrend(
		model.Between(allConsumption, newYear.Unix(), to),
		model.Between(allConsumption, newYear.AddDate(-1, 0, 0).Unix(), now.AddDate(-1, 0, 0).Unix()),
		model.Sum)
	twelveMonths := model.NewTrend(
		model.Between(allConsumption, now.AddDate(-1, 0, 0).Unix(), to),
		model.Between(allConsumption, now.AddDate(-2, 0, 0).Unix(), now.AddDate(-1, 0, 0).Unix()),
		model.Sum)

	return presenter.PresentHistory(dto.History{
		From:                  from,
		To:                    to,
		YearToDateTrend:       dto.NewTrend(yearToDate),
		TwelveMonthsTrend:     dto.NewTrend(twelveMonths),
		Consumption:           dto.NewSeasonalSeries("consumption", consumption),
		Temperature:           dto.NewSeasonalSeries(OUTSIDE_SERIES_LABEL, since(allTemperature, from)),
		MonthlyConsumption:    monthlyBuckets(model.GroupSum(allConsumption, model.ByMonth), from, now),
		YearlyConsumption:     dto.NewBuckets(model.GroupSum(consumption, model.ByYear)),
		HighConsumptionOnCold: dto.NewBuckets(highConsumptionOnCold),
		HotDays:               dto.NewBuckets(hotDays),
		ColdDayTemperature:    a.settings.ColdDayTemperature,
		HighConsumptionWatt:   a.settings.HighConsumptionWatt,
		HotDayTemperature:     a.settings.HotDayTemperature,
	})
}

// monthlyBuckets gives to every month of the period its season and its trend compared to the same month, a year before.
// The current month is not compared: it is not finished yet.
func monthlyBuckets(months []model.Bucket, from int64, now time.Time) []dto.Bucket {
	month := func(ts int64) time.Time {
		t := time.Unix(ts, 0).UTC()
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	}
	values := map[time.Time]float64{}
	for _, b := range months {
		values[month(b.TS)] = b.Value
	}

	buckets := []dto.Bucket{}
	for _, b := range months {
		if month(b.TS).Before(month(from)) {
			continue
		}
		bucket := dto.Bucket{Label: b.Label, Value: b.Value, Season: string(model.SeasonOf(b.TS))}
		if previous, ok := values[month(b.TS).AddDate(-1, 0, 0)]; ok && month(b.TS) != month(now.Unix()) {
			bucket.Trend = dto.NewTrend(&model.Trend{Value: b.Value, Previous: &previous})
		}
		buckets = append(buckets, bucket)
	}
	return buckets
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
