package port

import "github.com/sylvek/domotik/datastore/domain"

type StateRepository interface {
	Repository
	Store(domain.State) error
	Retrieve() (domain.State, error)
}
