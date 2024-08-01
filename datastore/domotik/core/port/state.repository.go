package port

import "github.com/sylvek/domotik/datastore/core/model"

type StateRepository interface {
	Repository
	Store(model.State) error
	Retrieve() (model.State, error)
}
