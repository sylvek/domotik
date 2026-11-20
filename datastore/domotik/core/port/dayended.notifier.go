package port

import "github.com/sylvek/domotik/datastore/core/dto"

// DayEndedNotifier is called when a day ends.
type DayEndedNotifier interface {
	Notify(dto.DayEnded) error
}
