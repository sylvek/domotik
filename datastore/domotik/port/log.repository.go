package port

import "github.com/sylvek/domotik/datastore/domain"

type LogRepository interface {
	Repository
	Store(domain.Log) error
}
