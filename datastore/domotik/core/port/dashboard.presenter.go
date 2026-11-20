package port

import "github.com/sylvek/domotik/datastore/core/dto"

type DashboardPresenter interface {
	PresentToday(dto.Today) error
	PresentQuarter(dto.Quarter) error
	PresentHistory(dto.History) error
}
