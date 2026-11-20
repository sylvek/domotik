package port

import "github.com/sylvek/domotik/datastore/core/model"

type LogPublisher interface {
	Publish(model.Log) error
}
