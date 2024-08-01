package port

import "github.com/sylvek/domotik/datastore/core/model"

type LogRepository interface {
	Repository
	Store(model.Log) error
}
