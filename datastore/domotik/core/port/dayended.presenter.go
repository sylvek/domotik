package port

import "github.com/sylvek/domotik/datastore/core/dto"

// DayEndedPresenter presents a day when it ends.
type DayEndedPresenter interface {
	Present(dto.DayEnded) error
}
