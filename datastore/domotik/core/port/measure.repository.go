package port

import "github.com/sylvek/domotik/datastore/core/model"

// MeasureRepository reads the stored measures.
// Timestamps are unix timestamps in seconds, `from` is included, `to` is excluded.
type MeasureRepository interface {
	// Last returns the most recent raw measure of a sensor.
	Last(name string, unit string, from int64, to int64) (*model.Point, error)
	// Series returns the raw measures by sensor name. An empty name means every sensor.
	Series(name string, unit string, from int64, to int64) (map[string][]model.Point, error)
	// Daily returns the daily aggregated measures of a sensor.
	Daily(name string, from int64, to int64) ([]model.Point, error)
	// DailySum returns, day by day, the sum of every aggregated measure having the given unit.
	DailySum(unit string, from int64, to int64) ([]model.Point, error)
}
