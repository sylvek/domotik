package port

import "github.com/sylvek/domotik/datastore/core/dto"

// DayEndedNotifier is called when a day ends, once by daily measure.
type DayEndedNotifier interface {
	Notify(dto.DayEnded) error
}
