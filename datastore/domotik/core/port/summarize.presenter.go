package port

import "github.com/sylvek/domotik/datastore/core/dto"

type SummarizePresenter interface {
	Present(dto.Output) error
}
