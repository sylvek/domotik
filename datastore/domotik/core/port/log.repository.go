package port

import "github.com/sylvek/domotik/datastore/core/model"

type LogRepository interface {
	Repository
	Store(model.Log) error
	// CloseDay aggregates the measures of the day that has just ended.
	// It returns false, and does nothing, while the current day is still the same.
	CloseDay() bool
}
