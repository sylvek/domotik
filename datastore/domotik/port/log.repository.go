package port

import "github.com/sylvek/domotik/datastore/domain/model"

type LogRepository interface {
	Repository
	Store(model.Log) error
}
