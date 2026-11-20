package dto

import "github.com/sylvek/domotik/datastore/core/model"

type Point struct {
	TS     int64   `json:"t"`
	Value  float64 `json:"v"`
	Season string  `json:"s,omitempty"`
}

// Trend compares a value with the value of the same period one year before.
// Previous and Change (0.1 means +10%) are null when the previous period is unknown.
type Trend struct {
	Value    float64  `json:"value"`
	Previous *float64 `json:"previous"`
	Change   *float64 `json:"change"`
}

type Series struct {
	Label  string  `json:"label"`
	Points []Point `json:"points"`
}

type Bucket struct {
	Label  string  `json:"label"`
	Value  float64 `json:"value"`
	Season string  `json:"season,omitempty"`
	// Trend compares the bucket with the same one, a year before.
	Trend *Trend `json:"trend,omitempty"`
}

type Reading struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

// Today describes the current day, from midnight to midnight (local time).
type Today struct {
	From                int64    `json:"from"`
	To                  int64    `json:"to"`
	LowTariff           *bool    `json:"lowTariff"`
	WattNow             *float64 `json:"wattNow"`
	WattConsumedToday   int64    `json:"wattConsumedToday"`
	EuroSpentToday      float64  `json:"euroSpentToday"`
	RatioLowTariffToday float64  `json:"ratioLowTariffToday"`
	// TypicalDay qualifies the consumption of the day, so far, by comparison with reference days.
	TypicalDay      string `json:"typicalDay"`
	TypicalDayLevel string `json:"typicalDayLevel"`
	// TypicalDayMaxWatt is the consumption (Wh) of the highest reference day.
	TypicalDayMaxWatt int64     `json:"typicalDayMaxWatt"`
	Inside            []Reading `json:"inside"`
	Outside           []Reading `json:"outside"`
	Temperatures      []Series  `json:"temperatures"`
	Power             []Series  `json:"power"`
}

// Quarter describes the last 90 or 180 days.
type Quarter struct {
	Days          int   `json:"days"`
	AvailableDays []int `json:"availableDays"`
	// trends compare the period with the same one, a year before.
	ConsumptionTrend *Trend   `json:"consumptionTrend"`
	CostTrend        *Trend   `json:"costTrend"`
	WaterTrend       *Trend   `json:"waterTrend"`
	TemperatureTrend *Trend   `json:"temperatureTrend"`
	From             int64    `json:"from"`
	To               int64    `json:"to"`
	Temperatures     []Series `json:"temperatures"`
	Consumption      Series   `json:"consumption"`
	Water            Series   `json:"water"`
	LowTariffRate    Series   `json:"lowTariffRate"`
	WeeklyCost       []Bucket `json:"weeklyCost"`
	WeeklyWaterCost  []Bucket `json:"weeklyWaterCost"`
	WaterEuroByLiter float64  `json:"waterEuroByLiter"`
}

// History describes the last 10 years.
type History struct {
	// YearToDateTrend compares the consumption since the 1st of january with last year's at the same date.
	YearToDateTrend *Trend `json:"yearToDateTrend"`
	// TwelveMonthsTrend compares the consumption of the last 12 months with the 12 months before.
	TwelveMonthsTrend     *Trend   `json:"twelveMonthsTrend"`
	From                  int64    `json:"from"`
	To                    int64    `json:"to"`
	Consumption           Series   `json:"consumption"`
	Temperature           Series   `json:"temperature"`
	MonthlyConsumption    []Bucket `json:"monthlyConsumption"`
	YearlyConsumption     []Bucket `json:"yearlyConsumption"`
	HighConsumptionOnCold []Bucket `json:"highConsumptionOnCold"`
	HotDays               []Bucket `json:"hotDays"`
	ColdDayTemperature    float64  `json:"coldDayTemperature"`
	HighConsumptionWatt   float64  `json:"highConsumptionWatt"`
	HotDayTemperature     float64  `json:"hotDayTemperature"`
}

func NewSeries(label string, points []model.Point) Series {
	series := Series{Label: label, Points: make([]Point, len(points))}
	for i, p := range points {
		series.Points[i] = Point{TS: p.TS, Value: p.Value}
	}
	return series
}

// NewSeasonalSeries gives its season to every point.
func NewSeasonalSeries(label string, points []model.Point) Series {
	series := NewSeries(label, points)
	for i, p := range points {
		series.Points[i].Season = string(model.SeasonOf(p.TS))
	}
	return series
}

// NewTrend returns nil when the trend is unknown.
func NewTrend(trend *model.Trend) *Trend {
	if trend == nil {
		return nil
	}
	return &Trend{Value: trend.Value, Previous: trend.Previous, Change: trend.Change()}
}

func NewBuckets(buckets []model.Bucket) []Bucket {
	result := make([]Bucket, len(buckets))
	for i, b := range buckets {
		result[i] = Bucket{Label: b.Label, Value: b.Value}
	}
	return result
}
