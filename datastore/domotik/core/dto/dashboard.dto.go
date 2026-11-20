package dto

import "github.com/sylvek/domotik/datastore/core/model"

type Point struct {
	TS    int64   `json:"t"`
	Value float64 `json:"v"`
}

type Series struct {
	Label  string  `json:"label"`
	Points []Point `json:"points"`
}

type Bucket struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

type Reading struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

// Today describes the current day, from midnight to midnight (local time).
type Today struct {
	From                int64     `json:"from"`
	To                  int64     `json:"to"`
	LowTariff           *bool     `json:"lowTariff"`
	WattNow             *float64  `json:"wattNow"`
	WattConsumedToday   int64     `json:"wattConsumedToday"`
	EuroSpentToday      float64   `json:"euroSpentToday"`
	RatioLowTariffToday float64   `json:"ratioLowTariffToday"`
	Inside              []Reading `json:"inside"`
	Outside             []Reading `json:"outside"`
	Temperatures        []Series  `json:"temperatures"`
	Power               []Series  `json:"power"`
}

// Quarter describes the last 90 days.
type Quarter struct {
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

// History describes the last 5 years.
type History struct {
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

func NewBuckets(buckets []model.Bucket) []Bucket {
	result := make([]Bucket, len(buckets))
	for i, b := range buckets {
		result[i] = Bucket{Label: b.Label, Value: b.Value}
	}
	return result
}
