package port

import "github.com/sylvek/domotik/datastore/core/dto"

// DayEndedNotifier is called when a day ends.
// The delivery (and its retries) is up to the implementation.
type DayEndedNotifier interface {
	Notify(dto.DayEnded) error
}
