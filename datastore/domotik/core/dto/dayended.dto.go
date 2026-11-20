package dto

import (
	"time"

	"github.com/sylvek/domotik/datastore/core/model"
)

// DayEnded describes a day that has just ended.
type DayEnded struct {
	// Timestamp is midnight UTC of the ended day, in seconds.
	Timestamp int64
	// WattConsumed is the electricity consumed during the day, in Wh.
	WattConsumed int64
}

// NewDayEndedFromState must be called before the state is reset for the new day.
// It returns nil when nothing was consumed during the day.
func NewDayEndedFromState(state model.State) *DayEnded {
	wattConsumed := state.DailySumHigh + state.DailySumLow
	if wattConsumed <= 0 || state.LastIndiceTS <= 0 {
		return nil
	}

	// the ended day is the day of the last index
	year, month, day := time.Unix(state.LastIndiceTS, 0).Date()

	return &DayEnded{
		Timestamp:    time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Unix(),
		WattConsumed: wattConsumed,
	}
}
